-- +goose Up
CREATE TABLE subscriptions (
    "id" BIGSERIAL PRIMARY KEY,
    "tenant_id" BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    "plan_id" BIGINT NOT NULL REFERENCES plans(id),

    -- Billing info
    "price" NUMERIC(10,2) NOT NULL,
    "original_price" NUMERIC(10,2) NOT NULL,
    "currency" VARCHAR(3) NOT NULL DEFAULT 'EGP',
    "billing_cycle" VARCHAR(50) NOT NULL DEFAULT 'monthly', -- monthly, yearly, one_time
    "auto_renew" BOOLEAN NOT NULL DEFAULT FALSE,

    -- Lifecycle
    "status" VARCHAR(20) NOT NULL DEFAULT 'active', -- active, expired, canceled, trial
    "start_date" TIMESTAMP NOT NULL DEFAULT NOW(),
    "end_date" TIMESTAMP,

    -- Snapshot of plan settings (so changes to plan don’t affect existing subs)
    "plan_limits" JSONB NOT NULL DEFAULT '{}'::jsonb,

    "expire_action" VARCHAR(50) NOT NULL DEFAULT 'block', -- enum-like
    "downgrade_to_plan" BIGINT REFERENCES plans(id), -- nullable
    "grace_period_days" INT NOT NULL DEFAULT 0,

    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS subscriptions CASCADE;
