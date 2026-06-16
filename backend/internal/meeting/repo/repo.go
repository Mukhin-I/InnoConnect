package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"gitlab.pg.innopolis.university/innoconnect-team/innoconnect/internal/meeting/entity"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

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
		return entity.Meeting{}, err
	}

	return meeting, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.Meeting, error) {
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
		ORDER BY meeting_time
	`

	rows, err := r.pool.Query(ctx, query)
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
