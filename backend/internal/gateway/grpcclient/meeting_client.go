package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/config"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/logger"
	meetingpb "gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/pb/meeting"
)

func NewMeetingClient() (meetingpb.MeetingServiceClient, error) {
	meeting_host := config.GetVar("GRPC_MEETING_CLIENT_HOST")
	meeting_port := config.GetVar("MEETING_SERVICE_PORT")
	meeting_url := meeting_host + ":" + meeting_port
	logger.Info("Creating gRPC client at " + meeting_url)
	conn, err := grpc.NewClient(
		meeting_url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return meetingpb.NewMeetingServiceClient(conn), nil
}