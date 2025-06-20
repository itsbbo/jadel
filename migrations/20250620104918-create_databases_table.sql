-- +migrate Up
create table if not exists databases (
    id bytea primary key not null,
    environment_id bytea not null references environments(id) on delete cascade,
    database_type varchar(100) not null,
    name varchar(255) not null,
    description varchar(255),
    image text not null,
    username varchar(255) not null,
    password text,
    port_mappings hstore,
    custom_config text,
    metadata jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists databases;