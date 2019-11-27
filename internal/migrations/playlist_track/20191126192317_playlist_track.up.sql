create table playlist_track (
    id bigserial not null primary key,
    playlist_id bigint not null references playlists(id) on delete cascade,
    track_id bigint not null references tracks(id) on delete cascade,
    created_at timestamp not null default now()
)
