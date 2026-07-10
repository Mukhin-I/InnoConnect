package transport

import (
	"context"
	"errors"

	"innoconnect/internal/chat/entity"
	pb "innoconnect/pkg/pb/chat"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ChatUsecase interface {
	GetOrCreateRequestChat(ctx context.Context, requestID, userID int64) (entity.Chat, error)
	GetMeetingChat(ctx context.Context, meetingID, userID int64) (entity.Chat, error)
	GetChat(ctx context.Context, chatID, userID int64) (entity.Chat, error)
	GetChats(ctx context.Context, userID int64) ([]entity.ChatPreview, error)
	GetMessages(ctx context.Context, chatID, userID int64) ([]entity.Message, error)
	SendMessage(ctx context.Context, chatID, userID int64, text string) (entity.Message, error)
	CreateMeetingChat(ctx context.Context, meetingID int64, chatName string, creator_id int64, creator_name string) (entity.Chat, error)
	AddToMeetingChat(ctx context.Context, meetingID int64, user_id int64, user_name string) (error)
}

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	usecase ChatUsecase
}

func NewChatServer(usecase ChatUsecase) *ChatServer {
	return &ChatServer{
		usecase: usecase,
	}
}

func (s *ChatServer) GetOrCreateRequestChat(
	ctx context.Context,
	req *pb.GetOrCreateRequestChatRequest,
) (*pb.ChatResponse, error) {

	chat, err := s.usecase.GetOrCreateRequestChat(
		ctx,
		req.GetRequestId(),
		req.GetUserId(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toChatResponse(chat), nil
}

func (s *ChatServer) CreateMeetingChat(
	ctx context.Context,
	req *pb.CreateMeetingChatRequest,
) (*pb.ChatResponse, error) {

	// 1. Validate input (basic safety)
	if req.GetMeetingId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "meeting_id is required")
	}

	if req.GetCreatorId() == 0 {
		return nil, status.Error(codes.Unauthenticated, "user_id is required")
	}

	// 2. Call usecase
	chat, err := s.usecase.CreateMeetingChat(
		ctx,
		req.GetMeetingId(),
		req.GetChatName(),
		req.GetCreatorId(),
		req.GetCreatorName(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 3. Map response
	return &pb.ChatResponse{
		ChatId: chat.ID,
		Type:   pb.ChatType_MEETING,
	}, nil
}

func (s *ChatServer) GetMeetingChat(
	ctx context.Context,
	req *pb.GetMeetingChatRequest,
) (*pb.ChatResponse, error) {

	chat, err := s.usecase.GetMeetingChat(
		ctx,
		req.GetMeetingId(),
		req.GetCreatorId(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toChatResponse(chat), nil
}

func (s *ChatServer) GetChat(
	ctx context.Context,
	req *pb.GetChatRequest,
) (*pb.GetChatResponse, error) {

	chat, err := s.usecase.GetChat(
		ctx,
		req.GetChatId(),
		req.GetUserId(),
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "chat not found")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toGetChatResponse(chat), nil
}

func (s *ChatServer) GetChats(
	ctx context.Context,
	req *pb.GetChatsRequest,
) (*pb.GetChatsResponse, error) {

	chats, err := s.usecase.GetChats(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetChatsResponse{
		Chats: make([]*pb.ChatPreview, 0, len(chats)),
	}

	for _, chat := range chats {
		response.Chats = append(response.Chats, toChatPreview(chat))
	}

	return response, nil
}

func (s *ChatServer) GetMessages(
	ctx context.Context,
	req *pb.GetMessagesRequest,
) (*pb.GetMessagesResponse, error) {

	messages, err := s.usecase.GetMessages(
		ctx,
		req.GetChatId(),
		req.GetUserId(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetMessagesResponse{
		Messages: make([]*pb.Message, 0, len(messages)),
	}

	for _, msg := range messages {
		response.Messages = append(response.Messages, toMessage(msg))
	}

	return response, nil
}

func (s *ChatServer) SendMessage(
	ctx context.Context,
	req *pb.SendMessageRequest,
) (*pb.Message, error) {

	message, err := s.usecase.SendMessage(
		ctx,
		req.GetChatId(),
		req.GetUserId(),
		req.GetText(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toMessage(message), nil
}

func (c *ChatServer) AddToMeetingChat(ctx context.Context, req *pb.CreateMeetingChatRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.usecase.AddToMeetingChat(ctx, req.MeetingId, req.CreatorId, req.CreatorName)
}