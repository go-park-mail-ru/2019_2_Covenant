create table subscriptions (
   id bigserial not null primary key,
   user_id bigint not null references users(id) on delete cascade,
   subscribed_to bigint not null references users(id) on delete cascade,
   created_at timestamp not null default now(),
   unique (user_id, subscribed_to)
)
