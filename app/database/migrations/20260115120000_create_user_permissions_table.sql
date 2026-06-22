-- +goose Up
CREATE TABLE "user_permissions" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL,
    "permission_key" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "fk_user_permissions_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE,
    CONSTRAINT "uk_user_permissions_user_key" UNIQUE ("user_id", "permission_key")
);

CREATE INDEX "idx_user_permissions_user_id" ON "user_permissions"("user_id");
CREATE INDEX "idx_user_permissions_key" ON "user_permissions"("permission_key");

-- +goose Down
DROP TABLE IF EXISTS "user_permissions" CASCADE;

