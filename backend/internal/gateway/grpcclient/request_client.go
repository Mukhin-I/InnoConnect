package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	requestpb "innoconnect/pkg/pb/request"
)

// NewRequestClient creates a new gRPC client for the request service
func NewRequestClient() (requestpb.RequestServiceClient, error) {
	requestHost := config.GetVar("GRPC_REQUEST_CLIENT_HOST")
	if requestHost == "" {
		requestHost = "request"
	}
	requestPort := config.GetVar("REQUEST_SERVICE_PORT")
	requestURL := requestHost + ":" + requestPort

	logger.Info("Creating gRPC client at " + requestURL)
	conn, err := grpc.NewClient(
		requestURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return requestpb.NewRequestServiceClient(conn), nil
}
