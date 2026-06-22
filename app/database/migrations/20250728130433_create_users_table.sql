-- +goose Up
CREATE TABLE "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "tenant_id" BIGINT, -- should be nullable cause we have a client with no tenant
    "client_id" BIGINT NULL,
    "client_type" VARCHAR(10) NULL,
    "name" VARCHAR(100) NOT NULL,
    "username" VARCHAR(255) UNIQUE,
    "email" VARCHAR(150) UNIQUE,
    "password" TEXT, -- should be nullable for social login
    "provider" VARCHAR(50), -- google, github, facebook, etc
    "provider_id" VARCHAR(255), -- user id from provider
    "avatar" TEXT, -- optional
    "email_verified_at" TIMESTAMP,
    "role"  VARCHAR(100) NOT NULL,
    "status"  SMALLINT NOT NULL DEFAULT 1,
    "deleted_at" TIMESTAMP DEFAULT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS "users" CASCADE;
