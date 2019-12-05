create table followers (
    id bigserial not null primary key,
    user_id bigint not null references users(id) on delete cascade,
    follower_id bigint not null references users(id) on delete cascade,
    created_at timestamp not null default now(),
    unique (user_id, follower_id)
)
