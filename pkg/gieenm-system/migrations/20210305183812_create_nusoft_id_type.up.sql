begin;

drop type if exists nusoft_id_type;
create type nusoft_id_type as (time varchar(15), id int);

commit;