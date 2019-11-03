create table users (
    id bigserial not null primary key,
    nickname varchar not null unique,
    email varchar not null unique,
    password varchar not null,
    avatar varchar not null default varchar '/resources/avatars/default.jpg',
    role int not null default 0,
    access int not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
