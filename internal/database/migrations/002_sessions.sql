-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE,
    user_ip TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

INSERT INTO sessions (id, user_id,user_ip,refresh_token, created_at, expires_at)
VALUES ('id_of_an_old_session','1716daab-5868-477e-9f51-0df2a0e925b7'::uuid ,'192.0.0.2','tokenHash',to_timestamp(1744813163),to_timestamp(1744813163));

INSERT INTO sessions (id, user_id,user_ip,refresh_token, created_at, expires_at)
VALUES ('6c13d4e5-ab60-40d6-81ef-5f51bd7d6f24','655b8db4-0a6e-4dd8-82c7-49112ec15a29'::uuid ,'192.0.0.2','8xW5aw+oD/sJvmUv6rg3kMViMvt5ZpQMxeQPPvnhM08=',to_timestamp('2025-04-17 12:07:18', 'YYYY-MM-DD HH:MI:SS'),to_timestamp('2025-04-18 12:07:18', 'YYYY-MM-DD HH:MI:SS'));

INSERT INTO sessions (id, user_id,user_ip,refresh_token, created_at, expires_at)
VALUES ('af9730be-973a-4685-bb02-a6bc9882ce6a','bdaff134-3b88-4596-86e2-4c55e106b77e'::uuid ,'144.20.0.1','eWLMfsIOZcthL4deoWp8lgB1YyxKyjjuv/4KJUtKfeU=',to_timestamp('2025-04-20 12:27:37', 'YYYY-MM-DD HH:MI:SS'),to_timestamp('2025-04-20 12:27:37', 'YYYY-MM-DD HH:MI:SS'));

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

COMMIT;