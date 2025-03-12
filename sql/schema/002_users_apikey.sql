-- +goose Up
alter table users add column api_key varchar(64) unique not null default (
    encode(sha256(random()::text::bytea), 'hex')
);

-- +goose Down
alter table users drop column api_key;
