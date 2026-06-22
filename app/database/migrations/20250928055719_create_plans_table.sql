-- +goose Up
CREATE TABLE plans (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL, 
    "features" JSONB,
    "is_active" BOOLEAN DEFAULT TRUE,
    "popular" BOOLEAN DEFAULT FALSE,
    "hidden" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS plans CASCADE;
