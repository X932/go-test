BEGIN;

alter table
    "user" drop column "password";

COMMIT;