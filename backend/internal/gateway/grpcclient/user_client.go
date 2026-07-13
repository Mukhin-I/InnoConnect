package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	userpb "innoconnect/pkg/pb/user"
)

func NewUserClient() (userpb.UserServiceClient, error) {
	userHost := config.GetVar("GRPC_USER_CLIENT_HOST")
	if userHost == "" {
		userHost = "user"
	}

	userPort := config.GetVar("USER_SERVICE_PORT")
	userURL := userHost + ":" + userPort

	logger.Info("Creating gRPC client at " + userURL)
	conn, err := grpc.NewClient(
		userURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return userpb.NewUserServiceClient(conn), nil
}
