-- +migrate Up
create table if not exists users (
    id bytea not null primary key,
    name varchar(255) not null,
    email varchar(255) not null unique,
    password text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists users;