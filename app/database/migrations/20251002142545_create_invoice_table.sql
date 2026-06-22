-- +goose Up
CREATE TABLE invoices (
    "id"              BIGSERIAL PRIMARY KEY,
    "tenant_id"       BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,

    "creator_id"      BIGINT NULL,
    "creator_type"     VARCHAR(50) NULL, -- system, user
    
    "user_id"         BIGINT NULL, -- the user who this invoice is for
    "user_type"       VARCHAR(50) NOT NULL,   -- e.g. "user"

    "order_id"        BIGINT REFERENCES orders(id) ON DELETE SET NULL,
    "transaction_id"  BIGINT REFERENCES transactions(id) ON DELETE SET NULL,

    "invoiceable_type" VARCHAR(50) NULL,  -- e.g. "Order", "Subscription", "Addon" if there is a related table for the action
    "invoiceable_id"   BIGINT NULL,

    -- Invoice Info
    "invoice_seq"     BIGINT NOT NULL,
    "invoice_number"  VARCHAR(50) UNIQUE NOT NULL,   -- e.g. "INV-{TYPE}-2025-0001"
    "type"            VARCHAR(50) NOT NULL,         -- "subscription", "upgrade", "downgrade", "renew_subscription", "renew_radius_user", "add_radius_user","manual_charge", "refund"
    "status"          VARCHAR(50) NOT NULL DEFAULT 'pending',  -- "pending", "paid", "partially_paid", "failed", "refunded", "canceled"

    "amount"          NUMERIC(10,2) NOT NULL,
    "paid"            NUMERIC(10,2) NOT NULL DEFAULT 0,
    "currency"        CHAR(3) NOT NULL DEFAULT 'EGP',

    -- Optional gateway info
    "due_date"        TIMESTAMP,
    "paid_at"         TIMESTAMP,

    "notes"           TEXT,
    "metadata"        JSONB,

    "created_at"      TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at"      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoices_tenant ON invoices(tenant_id);

-- +goose Down
DROP TABLE IF EXISTS invoices CASCADE;
