package transport

import meetingpb "gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/pb/meeting"

type Handler struct {
	meetingClient meetingpb.MeetingServiceClient
}

func NewHandler(client meetingpb.MeetingServiceClient) *Handler {
	return &Handler{
		meetingClient: client,
	}
}