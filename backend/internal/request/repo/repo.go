package repo

import (
	"context"

	"innoconnect/internal/request/entity"

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

func (r *Repository) Create(ctx context.Context, request entity.Request) (entity.Request, error) {
	query := `
		INSERT INTO requests (
			creator_id,
			creator_name,
			title,
			description,
			requester_address,
			type,
			deadline
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		request.CreatorID,
		request.CreatorName,
		request.Title,
		request.Description,
		request.RequesterAddress,
		request.Type,
		request.Deadline,
	).Scan(&request.ID)

	if err != nil {
		return entity.Request{}, err
	}

	return request, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.Request, error) {
	query := `
		SELECT
			id,
			creator_id,
			creator_name,
			title,
			description,
			requester_address,
			type,
			deadline
		FROM requests
		ORDER BY deadline
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []entity.Request

	for rows.Next() {
		var request entity.Request

		err := rows.Scan(
			&request.ID,
			&request.CreatorID,
			&request.CreatorName,
			&request.Title,
			&request.Description,
			&request.RequesterAddress,
			&request.Type,
			&request.Deadline,
		)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (entity.Request, error) {
	query := `
		SELECT
			id,
			creator_id,
			creator_name,
			title,
			description,
			requester_address,
			type,
			deadline
		FROM requests
		WHERE id = $1
	`

	var request entity.Request

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&request.ID,
		&request.CreatorID,
		&request.CreatorName,
		&request.Title,
		&request.Description,
		&request.RequesterAddress,
		&request.Type,
		&request.Deadline,
	)

	if err != nil {
		return entity.Request{}, err
	}

	return request, nil
}
