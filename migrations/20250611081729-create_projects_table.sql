-- +migrate Up
create table projects (
    id serial PRIMARY KEY not null,
    name varchar(255) not null,
    description text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists projects;