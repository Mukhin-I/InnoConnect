package transport

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"innoconnect/internal/meeting/entity"
	"innoconnect/internal/meeting/usecase"
	pb "innoconnect/pkg/pb/meeting"
)

type MeetingUsecase interface {
	Create(ctx context.Context, meeting entity.Meeting) (entity.Meeting, error)
	GetAll(ctx context.Context) ([]entity.Meeting, error)
	GetByID(ctx context.Context, id int64) (entity.Meeting, error)
	Delete(ctx context.Context, id int64, creatorID int64) error
	ApplyOnMeeting(ctx context.Context, userid int64, username string, id int64) error
}

type MeetingServer struct {
	pb.UnimplementedMeetingServiceServer
	usecase MeetingUsecase
}

func NewMeetingServer(usecase MeetingUsecase) *MeetingServer {
	return &MeetingServer{
		usecase: usecase,
	}
}

func (s *MeetingServer) CreateMeeting(ctx context.Context, req *pb.CreateMeetingRequest) (*pb.MeetingShort, error) {
	meetingTime, err := time.Parse(time.RFC3339, req.GetMeetingTime())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "meeting_time must be in RFC3339 format")
	}

	created, err := s.usecase.Create(ctx, entity.Meeting{
		CreatorID:   req.GetCreatorId(),
		CreatorName: req.GetCreatorName(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Type:        req.GetType(),
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		MeetingTime: meetingTime,
		MaxPeople:   req.MaxPeople,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toMeetingShort(created), nil
}

func (s *MeetingServer) GetMeetings(ctx context.Context, _ *pb.GetMeetingsRequest) (*pb.GetMeetingsResponse, error) {
	meetings, err := s.usecase.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GetMeetingsResponse{
		Meetings: make([]*pb.MeetingShort, 0, len(meetings)),
	}

	for _, meeting := range meetings {
		response.Meetings = append(response.Meetings, toMeetingShort(meeting))
	}

	return response, nil
}

func toMeetingShort(meeting entity.Meeting) *pb.MeetingShort {
	return &pb.MeetingShort{
		Id:        meeting.ID,
		Address:   meeting.Address,
		Type:      meeting.Type,
		Latitude:  meeting.Latitude,
		Longitude: meeting.Longitude,
	}
}

func (s *MeetingServer) GetMeeting(
	ctx context.Context,
	req *pb.GetMeetingRequest,
) (*pb.MeetingFull, error) {

	meeting, err := s.usecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.MeetingFull{
		Id:            meeting.ID,
		Title:         meeting.Title,
		Description:   meeting.Description,
		Type:          meeting.Type,
		Address:       meeting.Address,
		Latitude:      meeting.Latitude,
		Longitude:     meeting.Longitude,
		MeetingTime:   meeting.MeetingTime.Format(time.RFC3339),
		CurrentPeople: 0, // fill later
		MaxPeople:     meeting.MaxPeople,
		Creator: &pb.User{
			Id:   meeting.CreatorID,
			Name: meeting.CreatorName,
		},
		Participants: []*pb.User{},
	}, nil
}

func (s *MeetingServer) DeleteMeeting(ctx context.Context, req *pb.DeleteMeetingRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetCreatorId() == 0 {
		return nil, status.Error(codes.Unauthenticated, "creator_id is required")
	}

	err := s.usecase.Delete(ctx, req.GetId(), req.GetCreatorId())
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "meeting not found")
	}
	if errors.Is(err, usecase.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *MeetingServer) ApplyOnMeeting(ctx context.Context, req *pb.ApplyOnMeetingRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.usecase.ApplyOnMeeting(ctx, req.User.Id, req.User.Name, req.MeetingId)
}
