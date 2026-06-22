-- +goose Up
CREATE TABLE plan_prices (
    "id" BIGSERIAL PRIMARY KEY,
    "plan_id" BIGINT NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    "price" NUMERIC(10,2) NOT NULL DEFAULT 0.00, -- monthly fee
    "discount" NUMERIC(10,2) DEFAULT 0.00,          -- discount
    "currency" CHAR(15) NOT NULL,                  -- (USD, EUR, EGP)
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(plan_id, currency)                  -- prevent duplicate pricing
);

-- +goose Down
DROP TABLE IF EXISTS plan_prices CASCADE;