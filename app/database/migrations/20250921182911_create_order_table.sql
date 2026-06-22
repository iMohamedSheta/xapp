-- +goose Up
CREATE TABLE orders (
  "id" BIGSERIAL PRIMARY KEY,
  "tenant_id"       BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  "user_id"         BIGINT REFERENCES users(id) ON DELETE SET NULL,

  -- Polymorphic relation
  "orderable_type"  VARCHAR(100) NOT NULL,   -- "Ticket", "Subscription", "Addon"
  "orderable_id"    BIGINT NOT NULL,

  -- Order details
  "quantity"        INT NOT NULL DEFAULT 1,
  "unit_price"      NUMERIC(12,2) NOT NULL,
  "total_price"     NUMERIC(14,2) NOT NULL,
  "currency"        CHAR(3) NOT NULL DEFAULT 'EGP',
  "status"          SMALLINT NOT NULL DEFAULT 1, -- pending = 1, paid = 2, cancelled = 3

  -- Timestamps
  "created_at"      TIMESTAMP NOT NULL DEFAULT NOW(),
  "updated_at"      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Composite index for polymorphic lookups
CREATE INDEX idx_orders_orderable ON orders("orderable_type", "orderable_id");

-- +goose Down
DROP TABLE IF EXISTS orders CASCADE;
