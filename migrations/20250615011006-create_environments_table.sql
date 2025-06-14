-- +migrate Up
create table environments (
    id bytea PRIMARY KEY not null,
    name varchar(255) not null,
    project_id bytea not null references projects(id) on delete cascade,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists environments;