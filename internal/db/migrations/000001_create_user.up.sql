create table if not exists "user" (
    id int generated always as identity primary key,
    first_name text,
    last_name text,
    age int
);