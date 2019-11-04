create table tracks (
    id bigserial not null primary key,
    album_id bigint not null references albums(id),
    name varchar not null,
    duration time not null
);
