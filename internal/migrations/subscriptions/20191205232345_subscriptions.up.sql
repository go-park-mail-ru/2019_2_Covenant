create table subscriptions (
   id bigserial not null primary key,
   user_id bigint not null references users(id) on delete cascade,
   subscribed_to bigint not null references users(id) on delete cascade,
   created_at timestamp not null default now(),
   unique (user_id, subscribed_to),
   check (user_id != subscribed_to),
   constraint FK_SUBS_TO_SUBS FOREIGN KEY (subscribed_to) REFERENCES users(id),
   constraint FK_SUBS_TO_USERS FOREIGN KEY (user_id) REFERENCES users(id)
)
