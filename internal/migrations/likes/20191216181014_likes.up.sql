create table likes (
    id bigserial not null primary key,
    user_id bigint not null references users(id) on delete cascade,
    track_id bigint not null references tracks(id) on delete cascade,
    created_at timestamp not null default now(),
    unique (user_id, track_id)
)
