package app

import (
	"context"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"innoconnect/internal/meeting/repo"
	"innoconnect/internal/meeting/transport"
	"innoconnect/internal/meeting/usecase"
	"innoconnect/pkg/config"
	"innoconnect/pkg/logger"
	pb "innoconnect/pkg/pb/meeting"
	"google.golang.org/grpc"
)

func CreateServer() {
	meetingPort := config.GetVar("MEETING_SERVICE_PORT")
	dbPool, err := pgxpool.New(context.Background(), meetingDatabaseURL())
	if err != nil {
		logger.Error("Failed to connect to meeting database: " + err.Error())
		return
	}
	defer dbPool.Close()

	for i := 0; i < 10; i++ {
		err := dbPool.Ping(context.Background())
		if err == nil {
			logger.Info("Pinged db successfully")
			break
		}
		time.Sleep(2 * time.Second)
	}
	//if err := dbPool.Ping(context.Background()); err != nil {
		//logger.Error("Failed to ping meeting database: " + err.Error())
		//return
	//}

	meetingRepo := repo.New(dbPool)
	meetingUsecase := usecase.New(meetingRepo)
	meetingServer := transport.NewMeetingServer(meetingUsecase)

	listener, err := net.Listen("tcp", ":"+meetingPort)
	if err != nil {
		logger.Error("Failed to listen on meeting service port: " + err.Error())
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMeetingServiceServer(grpcServer, meetingServer)

	logger.Info("Meeting gRPC service running on :" + meetingPort)

	if err := grpcServer.Serve(listener); err != nil {
		logger.Error("Meeting gRPC service stopped: " + err.Error())
	}
}

func meetingDatabaseURL() string {
	host := config.GetVar("DB_MEETING_HOST")
	if host == "" {
		host = "db_meeting"
	}

	port := config.GetVar("DB_MEETING_PORT")
	user := config.GetVar("DB_MEETING_USER")
	password := config.GetVar("DB_MEETING_PASSWORD")

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   "meeting",
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}
