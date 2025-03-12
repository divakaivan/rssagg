-- +goose Up

insert into users (id, created_at, updated_at, name) 
values ('00000000-0000-0000-0000-000000000000', now(), now(), 'Default User');

-- +goose Down

delete from users where id = '00000000-0000-0000-0000-000000000000';
