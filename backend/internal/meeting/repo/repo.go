package repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"innoconnect/internal/meeting/entity"
	"innoconnect/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository struct for meeting database operations
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a new instance of the Repository with the provided database connection pool
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// Create inserts a new meeting into the database and returns the created meeting with its ID
func (r *Repository) Create(ctx context.Context, meeting entity.Meeting) (entity.Meeting, error) {
	query := `
		INSERT INTO meetings (
			creator_id,
			creator_name,
			title,
			description,
			type,
			address,
			latitude,
			longitude,
			meeting_time,
			max_people
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		meeting.CreatorID,
		meeting.CreatorName,
		meeting.Title,
		meeting.Description,
		meeting.Type,
		meeting.Address,
		meeting.Latitude,
		meeting.Longitude,
		meeting.MeetingTime,
		meeting.MaxPeople,
	).Scan(&meeting.ID)

	if err != nil {
		logger.Error(err.Error())
		return entity.Meeting{}, err
	}

	return meeting, nil
}

// GetAll retrieves all meetings from the database that are scheduled for the future
func (r *Repository) GetAll(ctx context.Context) ([]entity.Meeting, error) {
	now := time.Now()

	query := `
		SELECT
			id,
			creator_id,
			creator_name,
			title,
			description,
			type,
			address,
			latitude,
			longitude,
			meeting_time,
			max_people
		FROM meetings
		WHERE meeting_time > $1
		ORDER BY meeting_time
	`

	rows, err := r.pool.Query(ctx, query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []entity.Meeting

	for rows.Next() {
		var meeting entity.Meeting

		err := rows.Scan(
			&meeting.ID,
			&meeting.CreatorID,
			&meeting.CreatorName,
			&meeting.Title,
			&meeting.Description,
			&meeting.Type,
			&meeting.Address,
			&meeting.Latitude,
			&meeting.Longitude,
			&meeting.MeetingTime,
			&meeting.MaxPeople,
		)
		if err != nil {
			return nil, err
		}

		meetings = append(meetings, meeting)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return meetings, nil
}

// GetByID retrieves a meeting by its ID from the database
func (r *Repository) GetByID(ctx context.Context, id int64) (entity.Meeting, error) {
	query := `
		SELECT
			id,
			creator_id,
			creator_name,
			title,
			description,
			type,
			address,
			latitude,
			longitude,
			meeting_time,
			max_people
		FROM meetings
		WHERE id = $1
	`

	var meeting entity.Meeting

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&meeting.ID,
		&meeting.CreatorID,
		&meeting.CreatorName,
		&meeting.Title,
		&meeting.Description,
		&meeting.Type,
		&meeting.Address,
		&meeting.Latitude,
		&meeting.Longitude,
		&meeting.MeetingTime,
		&meeting.MaxPeople,
	)

	if err != nil {
		return entity.Meeting{}, err
	}

	return meeting, nil
}

// Update modifies an existing meeting in the database
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM meetings
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// ApplyOnMeeting adds a user as a participant to a meeting in the database
func (r *Repository) ApplyOnMeeting(ctx context.Context, userid int64, username string, id int64) error {
	// Check context first
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	query := `
        INSERT INTO meeting_participants 
        (meeting_id, user_id, user_name) 
        VALUES ($1, $2, $3)
    `

	_, err := r.pool.Exec(ctx, query, id, userid, username)
	if err != nil {
		// Check for duplicate entry
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("user already applied to this meeting")
		}
		return err
	}

	return nil
}
