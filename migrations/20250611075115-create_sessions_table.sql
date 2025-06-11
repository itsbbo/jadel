-- +migrate Up
create table sessions (
    id char(24) not null primary key,
    user_id bytea not null references users(id),
    ip_address varchar(45),
    user_agent text,
    expired_at timestamptz not null
);

-- +migrate Down
drop table if exists sessions;