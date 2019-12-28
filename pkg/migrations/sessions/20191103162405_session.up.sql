create table sessions (
    id bigserial not null primary key,
    user_id bigint not null references users(id) on delete cascade,
    expires timestamp not null default now() + interval '24 hours',
    data varchar not null,
    constraint FK_SESS_TO_USERS FOREIGN KEY (user_id) REFERENCES users(id)
)
