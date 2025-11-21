BEGIN;

alter table
    "user"
add
    column "password" text not null;

COMMIT;