create table playlists (
    id bigserial not null primary key,
    name varchar not null,
    description varchar,
    owner_id bigint not null references users(id) on delete cascade,
    photo varchar not null default varchar '/resources/photos/default_playlist.jpg',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
)
