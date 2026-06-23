package app

import (
	"context"
	"net"
	"net/url"
	"time"

	"innoconnect/internal/request/repo"
	"innoconnect/internal/request/transport"
	"innoconnect/internal/request/usecase"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	pb "innoconnect/pkg/pb/request"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func CreateServer() {
	requestPort := config.GetVar("REQUEST_SERVICE_PORT")
	dbPool, err := pgxpool.New(context.Background(), requestDatabaseURL())
	if err != nil {
		logger.Error("Failed to connect to request database: " + err.Error())
		return
	}
	defer dbPool.Close()

	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = dbPool.Ping(context.Background())
		if pingErr == nil {
			logger.Info("Pinged request db successfully")
			break
		}
		time.Sleep(2 * time.Second)
	}
	if pingErr != nil {
		logger.Error("Failed to ping request database: " + pingErr.Error())
		return
	}

	requestRepo := repo.New(dbPool)
	requestUsecase := usecase.New(requestRepo)
	requestServer := transport.NewRequestServer(requestUsecase)

	listener, err := net.Listen("tcp", ":"+requestPort)
	if err != nil {
		logger.Error("Failed to listen on request service port: " + err.Error())
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRequestServiceServer(grpcServer, requestServer)

	logger.Info("Request gRPC service running on :" + requestPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Error("Request gRPC service stopped: " + err.Error())
	}
}

func requestDatabaseURL() string {
	host := config.GetVar("DB_REQUEST_HOST")
	if host == "" {
		host = "db_request"
	}

	port := config.GetVar("DB_REQUEST_PORT")
	user := config.GetVar("DB_REQUEST_USER")
	password := config.GetVar("DB_REQUEST_PASSWORD")

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "request",
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}
