-- +goose Up
CREATE TABLE notifications (
    "id" BIGSERIAL PRIMARY KEY,
    "tenant_id" BIGINT, -- should be nullable cause we have a client with no tenant
    "notifiable_id" BIGINT NOT NULL,
    "notifiable_type" VARCHAR(255) NOT NULL,
    "type" VARCHAR(255) NOT NULL,
    "data" JSON NOT NULL,
    "read_at" TIMESTAMP NULL,
    "opened_at" TIMESTAMP NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS notifications CASCADE;
