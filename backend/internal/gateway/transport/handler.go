package transport

import meetingpb "innoconnect/pkg/pb/meeting"

type Handler struct {
	meetingClient meetingpb.MeetingServiceClient
}

func NewHandler(client meetingpb.MeetingServiceClient) *Handler {
	return &Handler{
		meetingClient: client,
	}
}