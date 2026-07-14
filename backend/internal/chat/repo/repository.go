package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"innoconnect/internal/chat/entity"
	"innoconnect/pkg/logger"
)

// Repository provides methods to interact with the chat database.
type Repository struct {
	db *pgxpool.Pool
}

// New creates a new instance of Repository with the provided database connection pool.
func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetChatByID retrieves a chat by its ID, including its participants.
func (r *Repository) GetChatByID(
	ctx context.Context,
	chatID int64,
) (entity.Chat, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var chat entity.Chat
	var chatType string

	err := r.db.QueryRow(ctx, `
		SELECT id, name, type
		FROM chats
		WHERE id = $1
	`, chatID).Scan(
		&chat.ID,
		&chat.Name,
		&chatType,
	)

	chat.Type = entity.ChatTypeFromString(chatType)

	if err != nil {
		return entity.Chat{}, err
	}

	participants, err := r.getParticipants(ctx, chatID)
	if err != nil {
		return entity.Chat{}, err
	}

	chat.Participants = participants
	return chat, nil
}

// GetChatsByUserID retrieves all chat previews for a given user ID, including the last message and participants for each chat.
func (r *Repository) GetChatsByUserID(
	ctx context.Context,
	userID int64,
) ([]entity.ChatPreview, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	logger.Info("Getting chats for a user " + strconv.FormatInt(userID, 10))
	rows, err := r.db.Query(ctx, `
		SELECT 
			c.id,
			c.name,
			c.type,
			c.related_id
		FROM chats c
		JOIN chat_participants p ON p.chat_id = c.id
		WHERE p.user_id = $1
		ORDER BY c.id DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.ChatPreview

	for rows.Next() {
		var chatID int64
		var name string
		var chatType string
		var relatedID int64

		if err := rows.Scan(&chatID, &name, &chatType, &relatedID); err != nil {
			return nil, err
		}

		participants, _ := r.getParticipants(ctx, chatID)
		lastMsg, _ := r.getLastMessage(ctx, chatID)

		logger.Info("Getting chatid: " + strconv.FormatInt(chatID, 10) + " Chatname " + name)
		result = append(result, entity.ChatPreview{
			ID:           chatID,
			Name:         name,
			Type:         entity.ChatTypeFromString(chatType),
			Participants: participants,
			LastMessage:  lastMsg,
		})
	}
	return result, nil
}

// GetMessages retrieves all messages for a given chat ID, ordered by the time they were sent.
func (r *Repository) GetMessages(
	ctx context.Context,
	chatID int64,
) ([]entity.Message, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT id, sender_id, sender_name, text, sent_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY sent_at ASC
	`, chatID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message

	for rows.Next() {
		var m entity.Message

		if err := rows.Scan(
			&m.ID,
			&m.Sender.ID,
			&m.Sender.Name,
			&m.Text,
			&m.SentAt,
		); err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

// SendMessage inserts a new message into the database for a given chat ID and sender ID, returning the created message.
func (r *Repository) SendMessage(
	ctx context.Context,
	chatID, senderID int64,
	text string,
) (entity.Message, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var msg entity.Message

	err := r.db.QueryRow(ctx, `
		INSERT INTO messages (chat_id, sender_id, sender_name, text)
		SELECT $1, $2, cp.user_name, $3
		FROM chat_participants cp
		WHERE cp.chat_id = $1 AND cp.user_id = $2
		RETURNING id, sender_id, sender_name, text, sent_at
	`, chatID, senderID, text).Scan(
		&msg.ID,
		&msg.Sender.ID,
		&msg.Sender.Name,
		&msg.Text,
		&msg.SentAt,
	)

	if err != nil {
		return entity.Message{}, err
	}

	return msg, nil
}

// GetRequestChat retrieves a chat of type 'REQUEST' associated with a specific request ID and user ID, including its participants.
func (r *Repository) GetRequestChat(
	ctx context.Context,
	requestID, userID int64,
) (entity.Chat, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var chat entity.Chat
	var chatType string

	err := r.db.QueryRow(ctx, `
		SELECT id, type
		FROM chats
		WHERE type = 'REQUEST'
		AND related_id = $1
	`, requestID).Scan(&chat.ID, &chatType)

	chat.Type = entity.ChatTypeFromString(chatType)

	if err != nil {
		return entity.Chat{}, err
	}

	chat.Participants, _ = r.getParticipants(ctx, chat.ID)
	return chat, nil
}

// CreateRequestChat creates a new chat of type 'REQUEST' associated with a specific request ID and user ID, and adds the user as a participant.
func (r *Repository) CreateRequestChat(
	ctx context.Context,
	requestID, userID int64, username string, chatName string,
) (entity.Chat, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Chat{}, err
	}
	defer tx.Rollback(ctx)

	var chatID int64

	err = tx.QueryRow(ctx, `
		INSERT INTO chats (type, name, related_id)
		VALUES ('REQUEST', $1, $2)
		RETURNING id
	`, chatName, requestID).Scan(&chatID)

	if err != nil {
		logger.Error(err.Error())
		return entity.Chat{}, err
	}

	err = r.AddParticipant(ctx, tx, chatID, userID, username)

	if err != nil {
		logger.Error(err.Error())
		return entity.Chat{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error(err.Error())
		return entity.Chat{}, err
	}

	return entity.Chat{
		ID:   chatID,
		Type: entity.RequestChat,
	}, nil
}

// getParticipants retrieves all participants for a given chat ID.
func (r *Repository) getParticipants(
	ctx context.Context,
	chatID int64,
) ([]entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, `
		SELECT user_id, user_name
		FROM chat_participants
		WHERE chat_id = $1
	`, chatID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// getLastMessage retrieves the last message for a given chat ID, if it exists.
func (r *Repository) getLastMessage(
	ctx context.Context,
	chatID int64,
) (*entity.Message, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var msg entity.Message

	err := r.db.QueryRow(ctx, `
		SELECT id, sender_id, sender_name, text, sent_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY sent_at DESC
		LIMIT 1
	`, chatID).Scan(
		&msg.ID,
		&msg.Sender.ID,
		&msg.Sender.Name,
		&msg.Text,
		&msg.SentAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &msg, nil
}

// GetChat retrieves a chat of a specific type associated with a given meeting ID, including its participants.
func (r *Repository) GetChat(
	ctx context.Context,
	relatedID int64,
	chatType string,
) (entity.Chat, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var chat entity.Chat
	err := r.db.QueryRow(ctx, `
		SELECT id
		FROM chats
		WHERE type = $1
		  AND related_id = $2
	`, chatType, relatedID).Scan(
		&chat.ID,
	)

	chat.Type = entity.ChatTypeFromString(chatType)

	if err != nil {
		logger.Error(err.Error())
		return entity.Chat{}, err
	}

	participants, err := r.getParticipants(ctx, chat.ID)
	if err != nil {
		logger.Error(err.Error())
		return entity.Chat{}, err
	}

	chat.Participants = participants

	return chat, nil
}

// CreateMeetingChat creates a new chat of type 'MEETING' associated with a specific meeting ID, adds the creator as a participant, and returns the created chat.
func (r *Repository) CreateMeetingChat(ctx context.Context, meetingID int64, chatName string, creatorID int64, creatorName string) (entity.Chat, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Chat{}, err
	}
	defer tx.Rollback(ctx)

	var chat entity.Chat
	var chatType string

	err = tx.QueryRow(ctx, `
		INSERT INTO chats (type, name, related_id)
		VALUES ('MEETING', $1, $2)
		RETURNING id, type, related_id
	`, chatName, meetingID).Scan(
		&chat.ID,
		&chatType,
		&chat.RelatedID,
	)

	chat.Type = entity.ChatTypeFromString(chatType)

	if err != nil {
		logger.Error(err.Error())
		return entity.Chat{}, err
	}

	// TODO: add name saving into gateway
	r.AddParticipant(ctx, tx, chat.ID, creatorID, creatorName)

	if err := tx.Commit(ctx); err != nil {
		logger.Error(err.Error())
		return entity.Chat{}, err
	}

	return chat, nil
}

// AddParticipant adds a user as a participant to a specific chat, optionally within an existing transaction.
func (r *Repository) AddParticipant(
	ctx context.Context,
	tx pgx.Tx,
	chatID int64,
	userID int64,
	userName string,
) error {
	newTransaction := false
	if tx == nil {
		newTransaction = true
		var err error
		tx, err = r.db.Begin(ctx)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		defer tx.Rollback(ctx)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO chat_participants (chat_id, user_id, user_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id, user_id) DO NOTHING
	`

	_, err := tx.Exec(ctx, stmt, chatID, userID, userName)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if newTransaction {
		if err := tx.Commit(ctx); err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	return nil
}

// RemoveParticipant removes a user from the participants of a specific chat.
func (r *Repository) GetParticipants(
	ctx context.Context,
	chatID int64,
) ([]entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.getParticipants(ctx, chatID)
}
