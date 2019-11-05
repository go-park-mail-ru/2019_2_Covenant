create table tracks (
    id bigserial not null primary key,
    album_id bigint not null references albums(id) on delete cascade,
    name varchar not null,
    duration time not null
);
