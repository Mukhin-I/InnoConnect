package transport

import (
<<<<<<< backend/internal/gateway/transport/handler.go
	chatpb "innoconnect/pkg/pb/chat"
	meetingpb "innoconnect/pkg/pb/meeting"
	requestpb "innoconnect/pkg/pb/request"
	userpb "innoconnect/pkg/pb/user"
)

type Handler struct {
	meetingClient meetingpb.MeetingServiceClient
	requestClient requestpb.RequestServiceClient
	chatClient    chatpb.ChatServiceClient
	userClient    userpb.UserServiceClient
}

func NewHandler(
	meetingClient meetingpb.MeetingServiceClient,
	requestClient requestpb.RequestServiceClient,
	chatClient chatpb.ChatServiceClient,
	userClient userpb.UserServiceClient,
) *Handler {
	return &Handler{
		meetingClient: meetingClient,
		requestClient: requestClient,
		chatClient:    chatClient,
		userClient:    userClient,
	}
}
