package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"innoconnect/internal/user/entity"
	"innoconnect/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
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
	repo      UserRepository
	jwtSecret string
}

func New(repo UserRepository, jwtSecret string) *Usecase {
	return &Usecase{
		repo:      repo,
		jwtSecret: jwtSecret,
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
	logger.Info("Generating jwt for a user " + strconv.FormatInt(user.ID, 10))
	token, err := u.generateToken(user.ID, user.Name)
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

func (u *Usecase) generateToken(userID int64, name string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"iat":     now.Unix(),
		"exp":     now.Add(time.Duration(tokenExpiresIn) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}
