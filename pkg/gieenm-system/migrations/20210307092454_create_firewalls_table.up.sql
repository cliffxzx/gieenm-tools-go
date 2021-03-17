begin;

create table if not exists firewalls (
  id serial primary key,
  uid char(32) not null unique default replace(cast(uuid_generate_v4() as char(36)), '-', ''),
  name varchar(36) not null unique,
  host varchar(128) not null,
  username varchar(64) not null,
  password varchar(64) not null,
  created_at timestamp not null default current_timestamp,
  modified_at timestamp not null default current_timestamp
);

CREATE TRIGGER
  update_modified_timestamp
BEFORE UPDATE ON
  firewalls
FOR EACH ROW EXECUTE PROCEDURE
  trigger_set_timestamp(modified_at)
;

commit;