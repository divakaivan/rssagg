-- name: CreateFeedFollow :one
insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetFeedFollows :many
select * from feed_follows where user_id = $1;

-- name: DeleteFeedFollow :exec
delete from feed_follows where id = $1 and user_id = $2;
