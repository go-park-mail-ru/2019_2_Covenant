create table playlist_track (
    id bigserial not null primary key,
    playlist_id bigint not null references playlists(id) on delete cascade,
    track_id bigint not null references tracks(id) on delete cascade,
    created_at timestamp not null default now(),
    unique (playlist_id, track_id),
    constraint FK_PT_TO_TRACKS FOREIGN KEY (track_id) REFERENCES tracks(id),
    constraint FK_PT_TO_PLAYLISTS FOREIGN KEY (playlist_id) REFERENCES playlists(id)
)
