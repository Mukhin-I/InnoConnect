package transport

import (
	"context"
	"time"

	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/internal/meeting/entity"
	pb "gitlab.pg.innopolis.university/innoconnect-team/innoconnect/pkg/pb/meeting"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MeetingUsecase interface {
	Create(ctx context.Context, meeting entity.Meeting) (entity.Meeting, error)
	GetAll(ctx context.Context) ([]entity.Meeting, error)
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
