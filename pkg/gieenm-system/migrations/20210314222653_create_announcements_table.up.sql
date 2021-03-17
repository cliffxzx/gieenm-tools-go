begin;

drop type if exists announce_level;
create type announce_level as enum ('INFO', 'NOTICE', 'WARNING', 'EMERGENCY');

create table if not exists announcements (
  id serial primary key,
  uid char(32) not null unique default replace(cast(uuid_generate_v4() as char(36)), '-', ''),
  title text not null,
  content text not null,
  announcer varchar(32) not null,
  level announce_level not null,
  created_at timestamp not null default current_timestamp,
  modified_at timestamp not null default current_timestamp
);

create trigger
  update_modified_timestamp
before update on
  announcements
for each row execute procedure
  trigger_set_timestamp(modified_at)
;

commit;