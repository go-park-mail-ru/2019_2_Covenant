\c postgres;
drop database if exists covenant_db;

create database covenant_db;

\c covenant_db

create table users (
    id bigserial not null primary key,
    nickname varchar not null unique,
    email varchar not null unique,
    password bytea not null,
    avatar varchar not null default varchar '/resources/avatars/default.jpg',
    role int not null default 0,
    access int not null default 0,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

CREATE UNIQUE INDEX nickname_unique_index on users (LOWER(nickname));

create table sessions (
    id bigserial not null primary key,
    user_id bigint not null references users(id) on delete cascade,
    expires timestamp not null default now() + interval '24 hours',
    data varchar not null
);

create table artists (
    id bigserial not null primary key,
    name varchar not null unique,
    photo varchar not null default varchar '/resources/photos/artists/default_artist.jpg'
);

create table albums (
    id bigserial not null primary key,
    artist_id bigint not null references artists(id) on delete cascade,
    name varchar not null,
    photo varchar not null default varchar '/resources/photos/albums/default_album.jpg',
    year date,
    unique (artist_id, name)
);

create table tracks (
    id bigserial not null primary key,
    album_id bigint not null references albums(id) on delete cascade,
    name varchar not null,
    duration time not null,
    path varchar not null default varchar '/resources/music/default.mp3',
    unique (album_id, name)
);

create table favourites (
    id bigserial not null primary key,
    user_id bigint not null references users(id) on delete cascade,
    track_id bigint not null references tracks(id) on delete cascade,
    created_at timestamp not null default now(),
    unique (user_id, track_id)
);

create table playlists (
    id bigserial not null primary key,
    name varchar not null,
    description varchar,
    owner_id bigint not null references users(id) on delete cascade,
    photo varchar not null default varchar '/resources/photos/playlists/default_playlist.jpg',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table playlist_track (
    id bigserial not null primary key,
    playlist_id bigint not null references playlists(id) on delete cascade,
    track_id bigint not null references tracks(id) on delete cascade,
    created_at timestamp not null default now(),
    unique (playlist_id, track_id)
);

create table subscriptions (
   id bigserial not null primary key,
   user_id bigint not null references users(id) on delete cascade,
   subscribed_to bigint not null references users(id) on delete cascade,
   created_at timestamp not null default now(),
   unique (user_id, subscribed_to),
   check (user_id != subscribed_to)
);

create table likes (
	    id bigserial not null primary key,
	    user_id bigint not null references users(id) on delete cascade,
	    track_id bigint not null references tracks(id) on delete cascade,
	    created_at timestamp not null default now(),
	    unique (user_id, track_id)
)