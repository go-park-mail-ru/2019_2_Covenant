create table artists (
    id bigserial not null primary key,
    name varchar not null unique,
    photo varchar not null default varchar '/resources/photos/artists/default_artist.jpg'
);
