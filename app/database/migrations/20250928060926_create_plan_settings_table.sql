-- +goose Up
CREATE TABLE plan_settings (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES plans(id) ON DELETE CASCADE,

    -- Limits
    plan_limits JSONB NOT NULL DEFAULT '{}'::jsonb,

    -- Expire action
    expire_action VARCHAR(20) NOT NULL DEFAULT 'block',  -- options: block, downgrade
    downgrade_to_plan BIGINT REFERENCES plans(id), -- if expire_action = downgrade

    -- Grace period
    grace_period_days INT DEFAULT 0,              -- days before doing expire_action

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS plan_settings CASCADE;
