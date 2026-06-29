package app

import (
	"context"
	"net"
	"net/url"
	"time"

	"innoconnect/internal/user/repo"
	"innoconnect/internal/user/transport"
	"innoconnect/internal/user/usecase"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	pb "innoconnect/pkg/pb/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func CreateServer() {
	userPort := config.GetVar("USER_SERVICE_PORT")
	dbPool, err := pgxpool.New(context.Background(), userDatabaseURL())
	if err != nil {
		logger.Error("Failed to connect to user database: " + err.Error())
		return
	}
	defer dbPool.Close()

	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = dbPool.Ping(context.Background())
		if pingErr == nil {
			logger.Info("Pinged user db successfully")
			break
		}
		time.Sleep(2 * time.Second)
	}
	if pingErr != nil {
		logger.Error("Failed to ping user database: " + pingErr.Error())
		return
	}

	userRepo := repo.New(dbPool)
	userUsecase := usecase.New(userRepo)
	userServer := transport.NewUserServer(userUsecase)

	listener, err := net.Listen("tcp", ":"+userPort)
	if err != nil {
		logger.Error("Failed to listen on user service port: " + err.Error())
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)

	logger.Info("User gRPC service running on :" + userPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Error("User gRPC service stopped: " + err.Error())
	}
}

func userDatabaseURL() string {
	host := config.GetVar("DB_USER_HOST")
	if host == "" {
		host = "db_user"
	}

	port := config.GetVar("DB_USER_PORT")
	user := config.GetVar("DB_USER_USER")
	password := config.GetVar("DB_USER_PASSWORD")

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "user",
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}
