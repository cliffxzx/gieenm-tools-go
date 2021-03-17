begin;

create table if not exists records (
  id serial primary key,
  uid char(32) not null unique default replace(cast(uuid_generate_v4() as char(36)), '-', ''),
  nusoft_id nusoft_id_type,
  user_id int not null references users(id) on update cascade,
  group_id int not null references record_groups(id) on update cascade,
  name varchar(36) not null,
  ip_addr inet not null,
  mac_addr macaddr not null,
  created_at timestamp not null default current_timestamp,
  modified_at timestamp not null default current_timestamp,
  unique(group_id, nusoft_id)
);

CREATE TRIGGER
  update_modified_timestamp
BEFORE UPDATE ON
  records
FOR EACH ROW EXECUTE PROCEDURE
  trigger_set_timestamp(modified_at)
;

commit;