create table albums (
    id bigserial not null primary key,
    artist_id bigint not null references artists(id) on delete cascade,
    name varchar not null,
    photo varchar not null default varchar '/resources/photos/albums/default_album.jpg',
    year date,
    unique (artist_id, name),
    constraint FK_ALBUMS_TO_ARTIST FOREIGN KEY (artist_id) REFERENCES artists(id)
);
