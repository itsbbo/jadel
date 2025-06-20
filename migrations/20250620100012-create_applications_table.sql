-- +migrate Up
create extension if not exists hstore;

create table if not exists applications (
    id bytea primary key not null,
    environment_id bytea not null references environments(id) on delete cascade,
    name varchar(255) not null,
    description varchar(255),
    build_tool varchar(100) not null,
    domain text not null,
    enable_https boolean not null default false,
    pre_deployment_script text,
    post_deployment_script text,
    port_mappings hstore,
    metadata jsonb,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

-- +migrate Down
drop table if exists applications;
drop extension if exists hstore;