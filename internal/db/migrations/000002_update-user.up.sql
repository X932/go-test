BEGIN;

alter table
    "user" drop column "age";

alter table
    "user"
add
    column "email" text not null;

alter table
    "user"
add
    constraint unique_user_email unique(email);

COMMIT;