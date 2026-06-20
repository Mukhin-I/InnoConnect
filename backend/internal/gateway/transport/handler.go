package transport

import (
requestpb "innoconnect/pkg/pb/request"
meetingpb "innoconnect/pkg/pb/meeting"
)

type Handler struct {
meetingClient meetingpb.MeetingServiceClient
requestClient requestpb.RequestServiceClient
}

func NewHandler(meetingClient meetingpb.MeetingServiceClient, requestClient requestpb.RequestServiceClient) *Handler {
return &Handler{
meetingClient: meetingClient,
requestClient: requestClient,
}
}
