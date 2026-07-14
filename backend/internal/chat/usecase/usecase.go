package usecase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"innoconnect/internal/chat/entity"
	"innoconnect/internal/chat/repo"
	"innoconnect/pkg/logger"
)

// ChatUsecase implements the business logic for chat operations
type ChatUsecase struct {
	repo *repository.Repository
}

// NewChatUsecase creates a new instance of ChatUsecase
func NewChatUsecase(repo *repository.Repository) *ChatUsecase {
	return &ChatUsecase{
		repo: repo,
	}
}

// GetChat retrieves a chat by its ID and user ID
func (u *ChatUsecase) GetChat(
	ctx context.Context,
	chatID,
	userID int64,
) (entity.Chat, error) {

	// Later maybe check that userID is a participant.
	return u.repo.GetChatByID(ctx, chatID)
}

// GetChats retrieves all chat previews for a given user ID
func (u *ChatUsecase) GetChats(
	ctx context.Context,
	userID int64,
) ([]entity.ChatPreview, error) {

	return u.repo.GetChatsByUserID(ctx, userID)
}

// GetMessages retrieves all messages for a given chat ID and user ID
func (u *ChatUsecase) GetMessages(
	ctx context.Context,
	chatID,
	userID int64,
) ([]entity.Message, error) {

	// Later: verify user belongs to the chat.
	return u.repo.GetMessages(ctx, chatID)
}

// SendMessage sends a message in a chat for a given user ID
func (u *ChatUsecase) SendMessage(
	ctx context.Context,
	chatID,
	userID int64,
	text string,
) (entity.Message, error) {

	// Later:
	// - validate text isn't empty
	// - verify membership

	return u.repo.SendMessage(ctx, chatID, userID, text)
}

// GetMeetingChat retrieves a meeting chat by its related meeting ID and user ID
func (u *ChatUsecase) GetMeetingChat(
	ctx context.Context,
	meetingID,
	userID int64,
) (entity.Chat, error) {

	// Later: verify user participates in meeting

	return u.repo.GetChat(ctx, meetingID, "MEETING")
}

// GetOrCreateRequestChat retrieves or creates a request chat by its related request ID and user ID
func (u *ChatUsecase) GetOrCreateRequestChat(
	ctx context.Context,
	requestID,
	userID int64,
	username string,
	chatName string,
) (entity.Chat, error) {

	chat, err := u.repo.GetRequestChat(ctx, requestID, userID)

	if err == nil {
		logger.Info("Get a request chat successfully")
		return chat, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return entity.Chat{}, err
	}
	logger.Info("Creating request chat")
	return u.repo.CreateRequestChat(ctx, requestID, userID, username, chatName)
}

// CreateMeetingChat creates a new meeting chat with the given details
func (u *ChatUsecase) CreateMeetingChat(ctx context.Context, meetingID int64, chatName string, creatorID int64, creatorName string) (entity.Chat, error) {

	chat, err := u.repo.GetChat(ctx, meetingID, "MEETING")
	if err == nil {
		return chat, nil
	}

	chat, err = u.repo.CreateMeetingChat(ctx, meetingID, chatName, creatorID, creatorName)
	if err != nil {
		return entity.Chat{}, err
	}

	return chat, nil
}

// AddToChat adds a user to an existing chat based on the related ID and chat type
func (u *ChatUsecase) AddToChat(ctx context.Context, relatedID int64, chatType string, user_id int64, user_name string) error {
	chat, err := u.repo.GetChat(ctx, relatedID, chatType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return u.repo.AddParticipant(ctx, nil, chat.ID, user_id, user_name)
}

// GetParticipants retrieves all participants for a given chat ID
func (u *ChatUsecase) GetParticipants(
	ctx context.Context,
	chatID int64,
) ([]entity.User, error) {

	return u.repo.GetParticipants(ctx, chatID)
}
