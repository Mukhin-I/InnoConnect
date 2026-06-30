package repo

import (
	"context"

	"innoconnect/internal/user/entity"

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
			name
		FROM users
		WHERE email = $1
	`

	var user entity.User

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
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
			name
		FROM users
		WHERE id = $1
	`

	var user entity.User

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
