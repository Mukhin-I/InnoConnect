package transport

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"innoconnect/internal/request/entity"
	pb "innoconnect/pkg/pb/request"
)

type RequestUsecase interface {
	Create(ctx context.Context, request entity.Request) (entity.Request, error)
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetByID(ctx context.Context, id int64) (entity.Request, error)
}

type RequestServer struct {
	pb.UnimplementedRequestServiceServer
	usecase RequestUsecase
}

func NewRequestServer(usecase RequestUsecase) *RequestServer {
	return &RequestServer{
		usecase: usecase,
	}
}

func (s *RequestServer) CreateRequest(ctx context.Context, req *pb.CreateRequestRequest) (*pb.RequestFull, error) {
	deadline, err := time.Parse(time.RFC3339, req.GetDeadline())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "deadline must be in RFC3339 format")
	}

	created, err := s.usecase.Create(ctx, entity.Request{
		CreatorID:        req.GetCreatorId(),
		CreatorName:      req.GetCreatorName(),
		Title:            req.GetTitle(),
		Description:      req.GetDescription(),
		RequesterAddress: req.GetRequesterAddress(),
		Type:             req.GetType(),
		Deadline:         deadline,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toRequestFull(created), nil
}

func (s *RequestServer) GetRequests(ctx context.Context, _ *pb.GetRequestsRequest) (*pb.GetRequestsResponse, error) {
	requests, err := s.usecase.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetRequestsResponse{
		Requests: make([]*pb.RequestShort, 0, len(requests)),
	}

	for _, request := range requests {
		response.Requests = append(response.Requests, toRequestShort(request))
	}

	return response, nil
}

func (s *RequestServer) GetRequest(ctx context.Context, req *pb.GetRequestRequest) (*pb.RequestFull, error) {
	request, err := s.usecase.GetByID(ctx, req.GetId())
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "request not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toRequestFull(request), nil
}

func toRequestShort(request entity.Request) *pb.RequestShort {
	return &pb.RequestShort{
		Id:       request.ID,
		Title:    request.Title,
		Type:     request.Type,
		Deadline: request.Deadline.Format(time.RFC3339),
	}
}

func toRequestFull(request entity.Request) *pb.RequestFull {
	return &pb.RequestFull{
		Id: request.ID,
		Creator: &pb.User{
			Id:   request.CreatorID,
			Name: request.CreatorName,
		},
		Title:            request.Title,
		Description:      request.Description,
		RequesterAddress: request.RequesterAddress,
		Type:             request.Type,
		Deadline:         request.Deadline.Format(time.RFC3339),
	}
}
