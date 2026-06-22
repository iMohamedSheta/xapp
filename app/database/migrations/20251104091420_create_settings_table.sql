-- +goose Up
CREATE TABLE "settings" (
  "id" BIGSERIAL PRIMARY KEY,
  "type" VARCHAR(255) NOT NULL,
  "model" VARCHAR(255),
  "model_id" BIGINT,
  "settings" JSONB DEFAULT '{}'::jsonb,
  "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT "settings_unique_type_model_model_id" UNIQUE ("type", "model", "model_id")
);

CREATE INDEX "settings_model_model_id_index" ON "settings" ("model", "model_id");

-- +goose Down
DROP TABLE IF EXISTS "settings" CASCADE;
