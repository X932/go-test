BEGIN;

alter table
    "user" drop constraint unique_user_email;

alter table
    "user" drop column "email";

alter table
    "user"
add
    column "age" int;

COMMIT;