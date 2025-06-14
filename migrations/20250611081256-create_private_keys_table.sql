-- +migrate Up
create table private_keys (
    id bytea PRIMARY KEY not null,
    name varchar(255) not null,
    user_id bytea not null references users(id) on delete cascade,
    description text,
    private_key text not null,
    is_git_related boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists private_keys;