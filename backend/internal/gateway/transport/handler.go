package transport

import (
	chatpb "innoconnect/pkg/pb/chat"
	meetingpb "innoconnect/pkg/pb/meeting"
	requestpb "innoconnect/pkg/pb/request"
)

type Handler struct {
	meetingClient meetingpb.MeetingServiceClient
	requestClient requestpb.RequestServiceClient
	chatClient    chatpb.ChatServiceClient
}

func NewHandler(
	meetingClient meetingpb.MeetingServiceClient,
	requestClient requestpb.RequestServiceClient,
	chatClient chatpb.ChatServiceClient,
) *Handler {
	return &Handler{
		meetingClient: meetingClient,
		requestClient: requestClient,
		chatClient:    chatClient,
	}
}
