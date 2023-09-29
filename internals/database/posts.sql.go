package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createPost = `
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description string
	PublishedAt time.Time
	FeedId      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedId,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedId,
	)
	return i, err
}

const getPostsByUser = `
SELECT * FROM posts WHERE feed_id = ANY($1) ORDER BY published_at DESC LIMIT $2
`

func (q *Queries) GetPostsByUser(ctx context.Context, feed_ids []uuid.UUID, limit int) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUser, pq.Array(feed_ids), limit)
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()
	posts := make([]Post, 0)
	for rows.Next() {
		var i Post
		err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedId,
		)
		if err != nil {
			return []Post{}, err
		}

		posts = append(posts, i)
	}
	return posts, nil
}
