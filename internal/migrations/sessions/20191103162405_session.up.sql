create table sessions (
    id bigserial not null primary key,
    user_id bigserial not null references users(id),
    expires timestamp not null default now() + interval '24 hours',
    data varchar not null
)
