-- +migrate Up
create table private_keys (
    id bytea PRIMARY KEY not null,
    name varchar(255) not null,
    user_id bytea not null references users(id) on delete cascade,
    description text,
    public_key text not null,
    private_key text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists private_keys;