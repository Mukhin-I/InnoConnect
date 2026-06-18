package usecase

import (
	"context"

	"innoconnect/internal/meeting/entity"
)

type MeetingRepository interface {
	Create(ctx context.Context, meeting entity.Meeting) (entity.Meeting, error)
	GetAll(ctx context.Context) ([]entity.Meeting, error)
	GetByID(ctx context.Context, id int64) (entity.Meeting, error)
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
