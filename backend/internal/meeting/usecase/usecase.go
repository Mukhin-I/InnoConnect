package usecase

import (
	"context"
	"errors"

	"innoconnect/internal/meeting/entity"
)

var ErrForbidden = errors.New("forbidden")

// Meeting repository interface
type MeetingRepository interface {
	Create(ctx context.Context, meeting entity.Meeting) (entity.Meeting, error)
	GetAll(ctx context.Context) ([]entity.Meeting, error)
	GetByID(ctx context.Context, id int64) (entity.Meeting, error)
	Delete(ctx context.Context, id int64) error
	ApplyOnMeeting(ctx context.Context, userid int64, username string, id int64) error
}

// Meeting usecase with business logic
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

// Delete checks if the user is the creator of the meeting before deleting it
func (u *Usecase) Delete(ctx context.Context, id int64, creatorID int64) error {
	meeting, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if meeting.CreatorID != creatorID {
		return ErrForbidden
	}

	return u.repo.Delete(ctx, id)
}

func (u *Usecase) ApplyOnMeeting(ctx context.Context, userid int64, username string, id int64) error {
	return u.repo.ApplyOnMeeting(ctx, userid, username, id)
}
