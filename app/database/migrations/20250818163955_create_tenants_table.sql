-- +goose Up
CREATE TABLE "tenants" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(150) NOT NULL,
    "status" SMALLINT NOT NULL DEFAULT 1,
    "balance" Numeric(10,2) NOT NULL DEFAULT 0.00,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS "tenants" CASCADE;
