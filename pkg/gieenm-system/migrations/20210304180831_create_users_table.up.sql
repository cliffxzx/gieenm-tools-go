begin;

drop type if exists role_enum;
create type role_enum as enum ('GUEST', 'STUDENT', 'TEACHER', 'MANAGER');

create table if not exists users (
  id serial primary key,
  uid char(32) not null unique default replace(cast(uuid_generate_v4() as char(36)), '-', ''),
  name varchar(36) not null,
  role role_enum not null,
  email varchar(64) unique not null,
  student_id varchar(64) unique,
  password varchar(64) not null,
  created_at timestamp not null default current_timestamp,
  modified_at timestamp not null default current_timestamp
);

CREATE TRIGGER
  update_modified_timestamp
BEFORE UPDATE ON
  users
FOR EACH ROW EXECUTE PROCEDURE
  trigger_set_timestamp(modified_at)
;

commit;