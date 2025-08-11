-- +goose Up
-- +goose StatementBegin
create table users (
    id serial primary key,
    name text not null,
    email text not null,
    password text not null,
    role text not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
