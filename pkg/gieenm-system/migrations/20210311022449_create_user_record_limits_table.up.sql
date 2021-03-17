begin;

create table if not exists user_record_limits (
  id serial primary key,
  uid char(32) not null unique default replace(cast(uuid_generate_v4() as char(36)), '-', ''),
  user_id int not null references users(id) on update cascade,
  group_id int not null references record_groups(id) on update cascade,
  record_max_count int not null,
  created_at timestamp not null default current_timestamp,
  modified_at timestamp not null default current_timestamp,
  unique(user_id, group_id)
);

create trigger
  update_modified_timestamp
before update on
  user_record_limits
for each row execute procedure
  trigger_set_timestamp(modified_at)
;

commit;