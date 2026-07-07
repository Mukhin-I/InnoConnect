package transport

import (
	"context"
	"errors"

	"innoconnect/internal/user/entity"
	"innoconnect/internal/user/usecase"
	"innoconnect/pkg/logger"
	pb "innoconnect/pkg/pb/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserUsecase interface {
	Register(ctx context.Context, email, password, name string) (entity.User, error)
	Login(ctx context.Context, email, password string) (usecase.LoginResult, error)
	GetCurrentUser(ctx context.Context, id int64) (entity.User, error)
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	usecase UserUsecase
}

func NewUserServer(usecase UserUsecase) *UserServer {
	return &UserServer{
		usecase: usecase,
	}
}

func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	_, err := s.usecase.Register(ctx, req.GetEmail(), req.GetPassword(), req.GetName())
	if isUniqueViolation(err) {
		return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		Message: "user registered successfully",
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := s.usecase.Login(ctx, req.GetEmail(), req.GetPassword())
	if errors.Is(err, usecase.ErrInvalidCredentials) {
		logger.Error(err.Error())
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}
	if err != nil {
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.LoginResponse{
		Token:     result.Token,
		Type:      result.Type,
		ExpiresIn: result.ExpiresIn,
	}

	return resp, nil
}

func (s *UserServer) GetCurrentUser(ctx context.Context, req *pb.GetCurrentUserRequest) (*pb.User, error) {
	user, err := s.usecase.GetCurrentUser(ctx, req.GetUserId())
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoUser(user), nil
}

func toProtoUser(user entity.User) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
