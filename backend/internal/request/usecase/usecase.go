package usecase

import (
	"context"

	"innoconnect/internal/request/entity"
)

type RequestRepository interface {
	Create(ctx context.Context, request entity.Request) (entity.Request, error)
	GetAll(ctx context.Context) ([]entity.Request, error)
	GetByID(ctx context.Context, id int64) (entity.Request, error)
}

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
