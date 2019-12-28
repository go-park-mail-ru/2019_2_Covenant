create table collections (
    id        bigserial not null primary key,
    name      varchar   not null,
    description varchar,
    photo     varchar   not null default varchar '/resources/photos/collections/default_collection.jpg',
    created_at timestamp not null default now()
);
