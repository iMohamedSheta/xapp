-- +goose Up
CREATE TABLE audit_logs (
    "id" BIGSERIAL PRIMARY KEY,
    "tenant_id" BIGINT NOT NULL,
    -- who did the action
    "user_id" BIGINT NULL,
    "user_type" VARCHAR(20) NULL,  -- 'user', 'radius_user', 'system'
    -- what object was affected
    "auditable_id" BIGINT NULL,
    "auditable_type" VARCHAR(50) NULL, -- 'radius_users', 'radius_cards', 'offers', 'nas', etc.
    "action" VARCHAR(30) NOT NULL,  -- 'create', 'update', 'delete', 'login', etc.
    "summary" VARCHAR(255) NOT NULL,
    "details" JSONB NOT NULL DEFAULT '{}'::jsonb,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS audit_logs CASCADE;
