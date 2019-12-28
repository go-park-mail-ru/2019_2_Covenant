create table tracks (
    id bigserial not null primary key,
    album_id bigint not null references albums(id) on delete cascade,
    name varchar not null,
    duration time not null default 0,
    path varchar not null default varchar '/resources/music/default.mp3',
    rating bigint not null default 0,
    unique (album_id, name),
    constraint FK_TRACKS_TO_ALBUM FOREIGN KEY (album_id) REFERENCES albums(id)
);
