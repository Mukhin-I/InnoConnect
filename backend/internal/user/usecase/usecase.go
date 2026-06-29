package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"innoconnect/internal/user/entity"

	"golang.org/x/crypto/bcrypt"
)

const (
	tokenType      = "Bearer"
	tokenExpiresIn = int32(3600)
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetByID(ctx context.Context, id int64) (entity.User, error)
}

type LoginResult struct {
	Token     string
	Type      string
	ExpiresIn int32
}

type Usecase struct {
	repo UserRepository
}

func New(repo UserRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Register(ctx context.Context, email, password, name string) (entity.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, err
	}

	return u.repo.Create(ctx, entity.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Name:         name,
	})
}

func (u *Usecase) Login(ctx context.Context, email, password string) (LoginResult, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := generateToken()
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Token:     token,
		Type:      tokenType,
		ExpiresIn: tokenExpiresIn,
	}, nil
}

func (u *Usecase) GetCurrentUser(ctx context.Context, id int64) (entity.User, error) {
	return u.repo.GetByID(ctx, id)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
