package repo

import (
	"context"

	"innoconnect/internal/user/entity"
	"innoconnect/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (
			email,
			password_hash,
			name
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.Name,
	).Scan(&user.ID)

	if err != nil {
		logger.Error(err.Error())
		return entity.User{}, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			name,
			created_requests_count,
			completed_requests_count
		FROM users
		WHERE email = $1
	`

	var user entity.User

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedRequestsCount,
		&user.CompletedRequestsCount,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (entity.User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			name,
			created_requests_count,
			completed_requests_count
		FROM users
		WHERE id = $1
	`

	var user entity.User

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedRequestsCount,
		&user.CompletedRequestsCount,
	)

	if err != nil {
		logger.Error(err.Error())
		return entity.User{}, err
	}

	return user, nil
}

func (r *Repository) IncrementCreatedRequestsCount(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET created_requests_count = created_requests_count + 1
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *Repository) IncrementCompletedRequestsCount(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET completed_requests_count = completed_requests_count + 1
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
