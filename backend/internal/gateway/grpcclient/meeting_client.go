package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/config"
	meetingpb "gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/pb/meeting"
)

func NewMeetingClient() (meetingpb.MeetingServiceClient, error) {
	conn, err := grpc.NewClient(
		"localhost:" + config.GetVar("GRPC_MEETING_CLIENT_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return meetingpb.NewMeetingServiceClient(conn), nil
}