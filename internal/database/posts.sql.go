// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
insert into posts (
    id, 
    created_at, 
    updated_at, 
    title, 
    summary, 
    content, 
    url, 
    published_at, 
    feed_id
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning id, created_at, updated_at, title, summary, content, url, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Summary     string
	Content     string
	Url         string
	PublishedAt string
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Summary,
		arg.Content,
		arg.Url,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Summary,
		&i.Content,
		&i.Url,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsForUser = `-- name: GetPostsForUser :many
select posts.id, posts.created_at, posts.updated_at, posts.title, posts.summary, posts.content, posts.url, posts.published_at, posts.feed_id from posts
join feed_follows on feed_follows.feed_id = posts.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc
limit $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Summary,
			&i.Content,
			&i.Url,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostsForUserWithPagination = `-- name: GetPostsForUserWithPagination :many
select posts.id, posts.created_at, posts.updated_at, posts.title, posts.summary, posts.content, posts.url, posts.published_at, posts.feed_id from posts
join feed_follows on feed_follows.feed_id = posts.feed_id
where feed_follows.user_id = $1
order by posts.published_at desc
limit $2
offset $3
`

type GetPostsForUserWithPaginationParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) GetPostsForUserWithPagination(ctx context.Context, arg GetPostsForUserWithPaginationParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUserWithPagination, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Summary,
			&i.Content,
			&i.Url,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
