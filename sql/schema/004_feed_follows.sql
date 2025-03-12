-- +goose Up

create table feed_follows (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null references users(id) on delete cascade,
    feed_id uuid not null references feeds(id) on delete cascade,
    unique(user_id, feed_id)
);

-- +goose Down

drop table feed_follows;
