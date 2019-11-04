create table albums (
    id bigserial not null primary key,
    artist_id bigint not null references artists(id),
    name varchar not null,
    photo varchar not null default varchar '/resources/photos/default_album.jpg',
    year date
);
