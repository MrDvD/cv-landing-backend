-- clear script
do $$
declare
  t text;
begin
  for t in
    select table_name
    from information_schema.tables
    where table_schema = 'public'
  loop
    execute 'drop table if exists ' || t || ' cascade';
  end loop;
end $$;

create type tag_type as enum(
  'core',
  'additional'
);

create type activity_type as enum(
  'project',
  'education',
  'event'
);

create table ACTIVITIES (
  id serial primary key,
  name text not null,
  subtitle text not null,
  description text,
  type activity_type not null,
  meta_label text,
  date_start date not null,
  date_end date
);

create table TAGS (
  name varchar(26) not null,
  type tag_type not null,
  activity_id int not null references ACTIVITIES(id) on delete cascade,
  unique (activity_id, name, type)
);