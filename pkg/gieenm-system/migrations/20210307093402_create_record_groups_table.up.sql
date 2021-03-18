begin;

create table if not exists record_groups (
  id serial primary key,
  uid char(32) not null unique default replace(cast(uuid_generate_v4() as char(36)), '-', ''),
  name varchar(36) not null,
  firewall_id int not null references firewalls on update cascade,
  nusoft_id nusoft_id_type not null,
  subnet inet not null,
  created_at timestamp not null default current_timestamp,
  modified_at timestamp not null default current_timestamp,
  unique(nusoft_id, firewall_id)
);

CREATE TRIGGER
  update_modified_timestamp
BEFORE UPDATE ON
  record_groups
FOR EACH ROW EXECUTE PROCEDURE
  trigger_set_timestamp(modified_at)
;

commit;