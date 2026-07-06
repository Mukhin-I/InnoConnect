package usecase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"innoconnect/internal/chat/entity"
	"innoconnect/internal/chat/repo"
)

type ChatUsecase struct {
	repo *repository.Repository
}

func NewChatUsecase(repo *repository.Repository) *ChatUsecase {
	return &ChatUsecase{
		repo: repo,
	}
}

func (u *ChatUsecase) GetChat(
	ctx context.Context,
	chatID,
	userID int64,
) (entity.Chat, error) {

	// Later you could check that userID is a participant.
	return u.repo.GetChatByID(ctx, chatID)
}

func (u *ChatUsecase) GetChats(
	ctx context.Context,
	userID int64,
) ([]entity.ChatPreview, error) {

	return u.repo.GetChatsByUserID(ctx, userID)
}

func (u *ChatUsecase) GetMessages(
	ctx context.Context,
	chatID,
	userID int64,
) ([]entity.Message, error) {

	// Later: verify user belongs to the chat.
	return u.repo.GetMessages(ctx, chatID)
}

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

func (u *ChatUsecase) GetMeetingChat(
	ctx context.Context,
	meetingID,
	userID int64,
) (entity.Chat, error) {

	// Later: verify user participates in meeting

	return u.repo.GetMeetingChat(ctx, meetingID)
}

func (u *ChatUsecase) GetOrCreateRequestChat(
	ctx context.Context,
	requestID,
	userID int64,
) (entity.Chat, error) {

	chat, err := u.repo.GetRequestChat(ctx, requestID, userID)

	if err == nil {
		return chat, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return entity.Chat{}, err
	}

	return u.repo.CreateRequestChat(ctx, requestID, userID)
}

func (u *ChatUsecase) CreateMeetingChat(ctx context.Context, meetingID, creatorID int64, creatorName string) (entity.Chat, error) {

	// 1. try find existing chat
	chat, err := u.repo.GetMeetingChat(ctx, meetingID)
	if err == nil {
		return chat, nil
	}

	// 2. create new chat
	chat, err = u.repo.CreateMeetingChat(ctx, meetingID, creatorID, creatorName)
	if err != nil {
		return entity.Chat{}, err
	}

	return chat, nil
}

func (u *ChatUsecase) AddToMeetingChat(ctx context.Context, meetingID int64, user_id int64, user_name string) (error) {
	chat, err := u.repo.GetMeetingChat(ctx, meetingID)
	if err != nil {
		return err
	}
	return u.repo.AddParticipant(ctx, nil, chat.ID, user_id, user_name)
}