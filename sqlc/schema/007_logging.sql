-- +goose Up

create table logs (
    id serial primary key,
    timestamp timestamp not null,
    caller_user_id uuid not null references users(id),
    method text not null,
    url text not null,
    status text not null,
    duration_ms bigint not null
);

-- +goose Down

drop table logs;
