-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE,
    user_ip TEXT CHECK (LENGTH(user_ip) <= 25) NOT NULL,
    refresh_token TEXT CHECK (LENGTH(user_ip) <= 50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NOT NULL
);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
BEGIN;

DROP TABLE sessions;

COMMIT;
-- +goose StatementEnd