package usecase

import (
	"context"
	"errors"

	"innoconnect/internal/request/entity"
)

var ErrForbidden = errors.New("forbidden")

type RequestRepository interface {
	Create(ctx context.Context, request entity.Request) (entity.Request, error)
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetByID(ctx context.Context, id int64) (entity.Request, error)
	Delete(ctx context.Context, id int64) error
	ApplyToRequest(ctx context.Context, requestID int64, userID int64, userName string) (string, error)
	CancelRequestApplication(ctx context.Context, requestID int64, userID int64, creatorID int64) error
}

// Usecase struct implements the business logic for requests
type Usecase struct {
	repo RequestRepository
}

func New(repo RequestRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Create(ctx context.Context, request entity.Request) (entity.Request, error) {
	return u.repo.Create(ctx, request)
}

func (u *Usecase) GetAll(ctx context.Context) ([]entity.Request, error) {
	return u.repo.GetAll(ctx)
}

func (u *Usecase) GetByID(ctx context.Context, id int64) (entity.Request, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *Usecase) Delete(ctx context.Context, id int64, creatorID int64) error {
	request, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if request.CreatorID != creatorID {
		return ErrForbidden
	}

	return u.repo.Delete(ctx, id)
}

func (u *Usecase) ApplyToRequest(
	ctx context.Context,
	requestID int64,
	userID int64,
	userName string,
) (string, error) {

	return u.repo.ApplyToRequest(ctx, requestID, userID, userName)
}

func (u *Usecase) CancelRequestApplication(
	ctx context.Context,
	requestID int64,
	userID int64,
	creatorID int64,
) error {

	return u.repo.CancelRequestApplication(ctx, requestID, userID, creatorID)
}
