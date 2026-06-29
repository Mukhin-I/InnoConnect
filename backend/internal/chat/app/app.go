package app

import (
	"context"
	"net"
	"net/url"
	"time"

	"innoconnect/internal/chat/repo"
	"innoconnect/internal/chat/transport"
	"innoconnect/internal/chat/usecase"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	pb "innoconnect/pkg/pb/chat"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func CreateServer() {
	port := config.GetVar("CHAT_SERVICE_PORT")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen: " + err.Error())
		return
	}

		dbPool, err := pgxpool.New(context.Background(), chatDatabaseURL())
	if err != nil {
		logger.Error("Failed to connect to chat database: " + err.Error())
		return
	}
	defer dbPool.Close()

		var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = dbPool.Ping(context.Background())
		if pingErr == nil {
			logger.Info("Pinged chat db successfully")
			break
		}
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		logger.Error("Failed to ping chat database: " + pingErr.Error())
		return
	}

	chatRepo := repository.New(dbPool)
	chatUsecase := usecase.NewChatUsecase(chatRepo)
	chatServer := transport.NewChatServer(chatUsecase)

		grpcServer := grpc.NewServer()

	pb.RegisterChatServiceServer(grpcServer, chatServer)

	logger.Info("Chat gRPC service running on :" + port)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve: " + err.Error())
	}
}

func chatDatabaseURL() string {
	host := config.GetVar("DB_CHAT_HOST")
	if host == "" {
		host = "db_chat"
	}

	port := config.GetVar("DB_CHAT_PORT")
	user := config.GetVar("DB_CHAT_USER")
	password := config.GetVar("DB_CHAT_PASSWORD")

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "chat",
	}

	q := databaseURL.Query()
	q.Set("sslmode", "disable")
	databaseURL.RawQuery = q.Encode()

	return databaseURL.String()
}