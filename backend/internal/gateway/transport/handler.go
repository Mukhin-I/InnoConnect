package transport

import (
    meetingpb "innoconnect/pkg/pb/meeting"
    requestpb "innoconnect/pkg/pb/request"
    userpb "innoconnect/pkg/pb/user"
)

type Handler struct {
    meetingClient meetingpb.MeetingServiceClient
    requestClient requestpb.RequestServiceClient
    userClient    userpb.UserServiceClient
}

func NewHandler(meetingClient meetingpb.MeetingServiceClient, requestClient requestpb.RequestServiceClient, userClient userpb.UserServiceClient) *Handler {
    return &Handler{
        meetingClient: meetingClient,
        requestClient: requestClient,
        userClient:    userClient,
    }
}
