-- +goose Up

-- Priority 1: invoices
CREATE INDEX idx_invoices_user_type_status ON invoices (user_id, user_type, status);
CREATE INDEX idx_invoices_status_created_at ON invoices (status, created_at DESC);
CREATE INDEX idx_invoices_order_id ON invoices (order_id);
CREATE INDEX idx_invoices_transaction_id ON invoices (transaction_id);
CREATE INDEX idx_invoices_tenant_created_at ON invoices (tenant_id, created_at);
CREATE INDEX idx_invoices_invoiceable ON invoices (invoiceable_type, invoiceable_id);

-- Priority 1: audit_logs
CREATE INDEX idx_audit_logs_tenant_created_at ON audit_logs (tenant_id, created_at DESC);
CREATE INDEX idx_audit_logs_auditable ON audit_logs (auditable_type, auditable_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_action_type ON audit_logs (action, auditable_type);

-- Priority 2: users
CREATE INDEX idx_users_tenant_id ON users (tenant_id);
CREATE INDEX idx_users_client ON users (client_id, client_type);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);
CREATE INDEX idx_users_role_active ON users (role) WHERE deleted_at IS NULL;

-- Priority 2: notifications
CREATE INDEX idx_notifications_notifiable ON notifications (notifiable_type, notifiable_id);
CREATE INDEX idx_notifications_unread ON notifications (notifiable_type, notifiable_id) WHERE read_at IS NULL;

-- Priority 2: orders
CREATE INDEX idx_orders_tenant_created_at ON orders (tenant_id, created_at DESC);
CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_status ON orders (status);

-- Priority 2: subscriptions
CREATE INDEX idx_subscriptions_tenant_id ON subscriptions (tenant_id);
CREATE INDEX idx_subscriptions_status ON subscriptions (status);

-- Priority 2: plan_settings / plans
CREATE INDEX idx_plan_settings_plan_id ON plan_settings (plan_id);
CREATE INDEX idx_plans_active ON plans (is_active, hidden) WHERE is_active = TRUE;

-- +goose Down
DROP INDEX IF EXISTS idx_plans_active;
DROP INDEX IF EXISTS idx_plan_settings_plan_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_tenant_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_tenant_created_at;
DROP INDEX IF EXISTS idx_notifications_unread;
DROP INDEX IF EXISTS idx_notifications_notifiable;
DROP INDEX IF EXISTS idx_users_role_active;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_client;
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_audit_logs_action_type;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_auditable;
DROP INDEX IF EXISTS idx_audit_logs_tenant_created_at;
DROP INDEX IF EXISTS idx_invoices_invoiceable;
DROP INDEX IF EXISTS idx_invoices_tenant_created_at;
DROP INDEX IF EXISTS idx_invoices_transaction_id;
DROP INDEX IF EXISTS idx_invoices_order_id;
DROP INDEX IF EXISTS idx_invoices_status_created_at;
DROP INDEX IF EXISTS idx_invoices_user_type_status;
