-- +migrate Up
create extension if not exists pgcrypto;
create extension if not exists pg_trgm;
create extension if not exists pg_stat_statements;

-- +migrate Down
drop extension if exists pg_stat_statements;
drop extension if exists pg_trgm;
drop extension if exists pgcrypto;