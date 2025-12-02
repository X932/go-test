BEGIN;

alter table
    article drop constraint if exists fk_owner;

drop table if exists article;

COMMIT;