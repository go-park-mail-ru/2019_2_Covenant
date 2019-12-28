create table collection_track (
    id            bigserial not null primary key,
    collection_id bigint    not null references collections (id) on delete cascade,
    track_id      bigint    not null references tracks (id) on delete cascade,
    created_at    timestamp not null default now(),
    unique (collection_id, track_id),
    constraint FK_CT_TO_TRACKS FOREIGN KEY (track_id) REFERENCES tracks (id),
    constraint FK_CT_TO_COLLECTIONS FOREIGN KEY (collection_id) REFERENCES collections (id)
)
