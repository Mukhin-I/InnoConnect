package grpcclient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	chatpb "innoconnect/pkg/pb/chat"
)

func NewChatClient() (chatpb.ChatServiceClient, error) {
	chatHost := config.GetVar("GRPC_CHAT_CLIENT_HOST")
	if chatHost == "" {
		chatHost = "chat"
	}

	chatPort := config.GetVar("CHAT_SERVICE_PORT")
	chatURL := chatHost + ":" + chatPort

	logger.Info("Creating gRPC client at " + chatURL)
	conn, err := grpc.NewClient(
		chatURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return chatpb.NewChatServiceClient(conn), nil
}
