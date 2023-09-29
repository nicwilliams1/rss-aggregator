package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// GET FEED BY ID
const getFeed = `-- name: GetFeed :one
SELECT * FROM feeds WHERE id = $1
`

func (q *Queries) GetFeed(ctx context.Context, id uuid.UUID) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeed, id)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserId,
		&i.LastFetchedAt,
	)
	return i, err
}

// GET FEEDS
const getFeeds = `-- name: GetFeeds :many
SELECT * FROM feeds ORDER BY created_at desc
`

func (q *Queries) GetFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return []Feed{}, err
	}
	defer rows.Close()
	feeds := make([]Feed, 0)
	for rows.Next() {
		var i Feed
		err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserId,
			&i.LastFetchedAt,
		)
		if err != nil {
			return []Feed{}, err
		}
		feeds = append(feeds, i)
	}
	return feeds, nil

}

// CREATE FEED
const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, name, url, user_id
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserId    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserId,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserId,
	)
	return i, err
}

// GET NEXT FEEDS TO FETCH
const getNextFeedsToFetch = `
SELECT * FROM feeds ORDER BY last_fetched_at NULLS FIRST LIMIT $1
`

func (q *Queries) GetNextFeedsToFetch(ctx context.Context, n int) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getNextFeedsToFetch, n)
	if err != nil {
		return []Feed{}, err
	}
	defer rows.Close()
	feeds := make([]Feed, 0)
	for rows.Next() {
		var i Feed
		err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserId,
			&i.LastFetchedAt,
		)
		if err != nil {
			return []Feed{}, err
		}
		feeds = append(feeds, i)
	}
	return feeds, nil
}

// MARK FEED AS FETCHED
const markFeedAsFetched = `
UPDATE feeds SET last_fetched_at = NOW(), updated_at = NOW() WHERE id = $1
`

func (q *Queries) MarkFeedAsFetched(ctx context.Context, feed_id uuid.UUID) error {
	sqlres, err := q.db.ExecContext(ctx, markFeedAsFetched, feed_id)
	if err != nil {
		return err
	}

	count, err := sqlres.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("update affected no rows")
	}

	return nil

}
