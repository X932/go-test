create table if not exists "article" (
    id serial primary key,
    title text,
    content text,
    tags jsonb,
    owner_id integer,
    constraint fk_owner foreign key (owner_id) references "user" (id) on delete
    set
        null
);