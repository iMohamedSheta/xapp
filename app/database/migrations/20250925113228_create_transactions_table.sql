-- +goose Up
CREATE TABLE transactions  (
  id                BIGSERIAL PRIMARY KEY,

  "tenant_id"       BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  "user_id"         BIGINT REFERENCES users(id) ON DELETE SET NULL,
  order_id          BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  gateway           VARCHAR(50) NOT NULL,      -- "paymob", "stripe", etc.

  -- Gateway info
  gateway_tx_id    VARCHAR(100),              -- Paymob transaction ID (id=347506848)
  gateway_order_id  VARCHAR(100),              -- Paymob order ID (order.id=390948532)
  merchant_order_id VARCHAR(100),              -- Your reference (e.g., "24")

  -- Amounts
  amount_cents      INT NOT NULL,
  currency          CHAR(3) NOT NULL DEFAULT 'EGP',
  captured_amount   INT DEFAULT 0,
  refunded_amount   INT DEFAULT 0,

  -- Status
  success           BOOLEAN DEFAULT false,
  status            VARCHAR(50),               -- "PAID", "FAILED", "REFUNDED"
  response_code     VARCHAR(50),               -- "APPROVED", "00", etc.
  message           VARCHAR(255),              -- Gateway message

  -- Payment Method Info (non-sensitive)
  payment_method    VARCHAR(50),               -- "card", "wallet", "fawry"
  card_type         VARCHAR(50),               -- "MASTERCARD"
  card_last4        VARCHAR(10),               -- "2346"

  -- Gateway Metadata
  authorize_id      VARCHAR(50),               -- "038788"
  receipt_no        VARCHAR(50),               -- "526804038788"
  batch_no          VARCHAR(20),               -- "20250925"
  gateway_ref_id    VARCHAR(100),              -- "123456789" (migs_transaction.transactionId)

  -- Customer Snapshot (optional)
  customer_email    VARCHAR(100),
  customer_phone    VARCHAR(20),
  customer_name     VARCHAR(100),

  -- Full response for debugging
  raw_response      JSONB,

  -- Timestamps
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_order_id ON transactions(order_id);
CREATE INDEX idx_transactions_gateway_tx ON transactions(gateway, gateway_tx_id);

-- +goose Down
DROP TABLE IF EXISTS transactions CASCADE;