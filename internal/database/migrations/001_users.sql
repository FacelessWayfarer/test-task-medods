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



CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT CHECK (LENGTH(email) <= 50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);

COMMIT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
BEGIN;

DROP TABLE users;

COMMIT;
-- +goose StatementEnd

