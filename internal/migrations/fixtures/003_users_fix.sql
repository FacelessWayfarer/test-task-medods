-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

BEGIN;

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;
SET log_error_verbosity = VERBOSE;

INSERT INTO users (id, email, created_at, updated_at)
VALUES ('1716daab-5868-477e-9f51-0df2a0e925b7'::uuid,'sjy1nmkmals2ldia9wnnla@gmail.com',to_timestamp(1744813161),to_timestamp(1744813161));


INSERT INTO users (id, email, created_at, updated_at)
VALUES ('655b8db4-0a6e-4dd8-82c7-49112ec15a29'::uuid,'252ynq3o5qozmawnnla@mail.ru',to_timestamp(1744815890),to_timestamp(1744815890));

INSERT INTO users (id, email, created_at, updated_at)
VALUES ('bdaff134-3b88-4596-86e2-4c55e106b77e'::uuid,'kasuuaklqozmawnnla@mail.ru',to_timestamp(1744815896),to_timestamp(1744815897));

COMMIT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
BEGIN;

DELETE FROM users WHERE id = '1716daab-5868-477e-9f51-0df2a0e925b7';

DELETE FROM users WHERE id = '655b8db4-0a6e-4dd8-82c7-49112ec15a29';

DELETE FROM users WHERE id = 'bdaff134-3b88-4596-86e2-4c55e106b77e';

COMMIT;
-- +goose StatementEnd

