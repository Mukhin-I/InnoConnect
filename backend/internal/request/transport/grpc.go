package transport

import (
	"context"
	"errors"
	"time"

	"innoconnect/internal/request/entity"
	"innoconnect/internal/request/usecase"
	"innoconnect/pkg/logger"
	pb "innoconnect/pkg/pb/request"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// RequestUsecase defines the interface for request business logic
type RequestUsecase interface {
	Create(ctx context.Context, request entity.Request) (entity.Request, error)
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetByID(ctx context.Context, id int64) (entity.Request, error)
	Delete(ctx context.Context, id int64, creatorID int64) error
	ApplyToRequest(ctx context.Context, requestID int64, userID int64, userName string) (string, error)
	CancelRequestApplication(ctx context.Context, requestID int64, userID int64, creatorID int64) error
}

// RequestServer implements the gRPC server for the request service
type RequestServer struct {
	pb.UnimplementedRequestServiceServer
	usecase RequestUsecase
}

// NewRequestServer creates a new instance of RequestServer with the provided usecase
func NewRequestServer(usecase RequestUsecase) *RequestServer {
	return &RequestServer{
		usecase: usecase,
	}
}

// CreateRequest handles the creation of a new request, validating the deadline format and returning the created request
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

// GetRequests retrieves all requests from the database that have a deadline in the future
func (s *RequestServer) GetRequests(ctx context.Context, _ *pb.GetRequestsRequest) (*pb.GetRequestsResponse, error) {
	requests, err := s.usecase.GetAll(ctx)
	if err != nil {
		logger.Error(err.Error())
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

// GetRequest retrieves a request by its ID from the database
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

// DeleteRequest checks if the user is the creator of the request before deleting it
func (s *RequestServer) DeleteRequest(ctx context.Context, req *pb.DeleteRequestRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetCreatorId() == 0 {
		return nil, status.Error(codes.Unauthenticated, "creator_id is required")
	}

	err := s.usecase.Delete(ctx, req.GetId(), req.GetCreatorId())
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "request not found")
	}
	if errors.Is(err, usecase.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// toRequestShort converts an entity.Request to a pb.RequestShort
func toRequestShort(request entity.Request) *pb.RequestShort {
	return &pb.RequestShort{
		Id:        request.ID,
		CreatorId: request.CreatorID,
		Title:     request.Title,
		Type:      request.Type,
		Deadline:  request.Deadline.Format(time.RFC3339),
	}
}

// toRequestFull converts an entity.Request to a pb.RequestFull
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
		Status:           request.Status,
		Deadline:         request.Deadline.Format(time.RFC3339),
	}
}

// ApplyToRequest allows a user to apply to a request, changing its status to "IN PROGRESS"
func (s *RequestServer) ApplyToRequest(
	ctx context.Context,
	req *pb.ApplyToRequestRequest,
) (*pb.ApplyToRequestResponse, error) {

	if req.GetRequestId() == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"request_id is required",
		)
	}

	if req.GetUserId() == 0 {
		return nil, status.Error(
			codes.Unauthenticated,
			"user_id is required",
		)
	}

	req_title, err := s.usecase.ApplyToRequest(
		ctx,
		req.GetRequestId(),
		req.GetUserId(),
		req.GetUserName(),
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(
				codes.NotFound,
				"request not found",
			)
		}

		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	return &pb.ApplyToRequestResponse{
		ReqTitle: req_title,
	}, nil
}

// CancelRequestApplication allows a user to cancel their application to a request, changing its status back to "PENDING"
func (s *RequestServer) CancelRequestApplication(
	ctx context.Context,
	req *pb.CancelRequestApplicationRequest,
) (*emptypb.Empty, error) {

	if req.GetRequestId() == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"request_id is required",
		)
	}

	if req.GetUserId() == 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"user_id is required",
		)
	}

	if req.GetCreatorId() == 0 {
		return nil, status.Error(
			codes.Unauthenticated,
			"creator_id is required",
		)
	}

	err := s.usecase.CancelRequestApplication(
		ctx,
		req.GetRequestId(),
		req.GetUserId(),
		req.GetCreatorId(),
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(
				codes.NotFound,
				"application not found",
			)
		}

		return nil, status.Error(
			codes.Internal,
			err.Error(),
		)
	}

	return &emptypb.Empty{}, nil
}
