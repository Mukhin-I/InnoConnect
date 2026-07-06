package usecase

import (
	"context"

	"innoconnect/internal/meeting/entity"
)

type MeetingRepository interface {
	Create(ctx context.Context, meeting entity.Meeting) (entity.Meeting, error)
	GetAll(ctx context.Context) ([]entity.Meeting, error)
	GetByID(ctx context.Context, id int64) (entity.Meeting, error)
	ApplyOnMeeting(ctx context.Context, userid int64, username string, id int64) error
}

type Usecase struct {
	repo MeetingRepository
}

func New(repo MeetingRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Create(ctx context.Context, meeting entity.Meeting) (entity.Meeting, error) {
	return u.repo.Create(ctx, meeting)
}

func (u *Usecase) GetAll(ctx context.Context) ([]entity.Meeting, error) {
	return u.repo.GetAll(ctx)
}

func (u *Usecase) GetByID(ctx context.Context, id int64) (entity.Meeting, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *Usecase) ApplyOnMeeting(ctx context.Context, userid int64, username string, id int64) error {
	return u.repo.ApplyOnMeeting(ctx, userid, username, id)
}
