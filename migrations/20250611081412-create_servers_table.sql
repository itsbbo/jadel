-- +migrate Up
create table servers (
    id bytea PRIMARY KEY not null,
    name varchar(255) not null,
    user_id bytea not null references users(id) on delete cascade,
    description text,
    ip inet not null,
    port integer not null default 22,
    username varchar(255) not null default 'root',
    private_key_id bytea not null references private_keys(id),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists servers;