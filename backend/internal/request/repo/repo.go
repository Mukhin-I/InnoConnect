package repo

import (
	"context"

	"innoconnect/internal/request/entity"
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
		logger.Error(err.Error())
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
			logger.Error(err.Error())
			return nil, err
		}

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		logger.Error(err.Error())
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
			status,
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
		&request.Status,
		&request.Deadline,
	)

	if err != nil {
		return entity.Request{}, err
	}

	return request, nil
}

func (r *Repository) ApplyToRequest(
	ctx context.Context,
	requestID int64,
	userID int64,
	userName string,
) (string, error) {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)


	// Add user application
	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO request_applications (
			request_id,
			user_id,
			user_name
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (request_id, user_id)
		DO NOTHING
		`,
		requestID,
		userID,
		userName,
	)

	if err != nil {
		return "", err
	}


	// Change request status
	var title string
	err = tx.QueryRow(
		ctx,
		`
		UPDATE requests
		SET status = 'IN PROGRESS'
		WHERE id = $1
		RETURNING title
		`,
		requestID,
	).Scan(&title)

	if err != nil {
		return "", err
	}

	return title, tx.Commit(ctx)
}

func (r *Repository) CancelRequestApplication(
	ctx context.Context,
	requestID int64,
	userID int64,
	creatorID int64,
) error {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)


	// Cancel application only if creator owns request
	result, err := tx.Exec(
		ctx,
		`
		UPDATE request_applications
		SET status = 'CANCELLED'
		WHERE request_id = $1
		  AND user_id = $2
		  AND EXISTS (
				SELECT 1
				FROM requests
				WHERE id = $1
				  AND creator_id = $3
		  )
		`,
		requestID,
		userID,
		creatorID,
	)

	if err != nil {
		return err
	}


	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}


	// Return request back to pending
	_, err = tx.Exec(
		ctx,
		`
		UPDATE requests
		SET status = 'PENDING'
		WHERE id = $1
		`,
		requestID,
	)

	if err != nil {
		return err
	}


	return tx.Commit(ctx)
}