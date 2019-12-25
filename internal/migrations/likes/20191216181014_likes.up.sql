create table likes (
    id bigserial not null primary key,
    user_id bigint not null references users(id) on delete cascade,
    track_id bigint not null references tracks(id) on delete cascade,
    created_at timestamp not null default now(),
    unique (user_id, track_id),
    constraint FK_LIKES_TO_TRACKS FOREIGN KEY (track_id) REFERENCES tracks(id),
    constraint FK_LIKES_TO_USERS FOREIGN KEY (user_id) REFERENCES users(id)
)
