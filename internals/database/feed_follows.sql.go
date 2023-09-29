package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// CREATE FEED FOLLOW
const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uuid.UUID
	FeedId    uuid.UUID
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserId,
		arg.FeedId,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserId,
		&i.FeedId,
	)
	return i, err
}

// DELETE FEED FOLLOW
const deleteFeedFollow = `-- name: DeleteFeedFollow :one
DELETE FROM feed_follows WHERE id = $1
`

func (q *Queries) DeleteFeedFollow(ctx context.Context, feed_id uuid.UUID) error {
	sqlres, err := q.db.ExecContext(ctx, deleteFeedFollow, feed_id)
	if err != nil {
		return err
	}

	count, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("delete affected no rows")
	}

	return nil
}

const getFeedFollowsByUserId = `-- name: GetFeedFollowsByUserId :many
SELECT * from feed_follows WHERE user_id = $1
`

func (q *Queries) GetFeedFollowsByUserId(ctx context.Context, user_id uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUserId, user_id)
	if err != nil {
		return []FeedFollow{}, err
	}

	defer rows.Close()
	feedFollows := make([]FeedFollow, 0)
	for rows.Next() {
		var i FeedFollow
		err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserId,
			&i.FeedId,
		)
		if err != nil {
			return []FeedFollow{}, err
		}
		feedFollows = append(feedFollows, i)
	}
	return feedFollows, nil

}
