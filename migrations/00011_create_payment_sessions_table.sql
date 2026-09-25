-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN IF NOT EXISTS payment_status VARCHAR(50) NOT NULL DEFAULT 'UNPAID',
ADD COLUMN IF NOT EXISTS stripe_session_id VARCHAR(255);

CREATE TABLE IF NOT EXISTS payment_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_session_id VARCHAR(255) NOT NULL UNIQUE,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    table_id UUID NOT NULL REFERENCES tables(id) ON DELETE RESTRICT,
    order_payload TEXT NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(10) NOT NULL DEFAULT 'aud',
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID
);

CREATE INDEX IF NOT EXISTS idx_payment_sessions_stripe_session_id ON payment_sessions(stripe_session_id);
CREATE INDEX IF NOT EXISTS idx_payment_sessions_restaurant_id ON payment_sessions(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_payment_sessions_table_id ON payment_sessions(table_id);
CREATE INDEX IF NOT EXISTS idx_payment_sessions_status ON payment_sessions(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_sessions;
ALTER TABLE orders DROP COLUMN IF EXISTS stripe_session_id;
ALTER TABLE orders DROP COLUMN IF EXISTS payment_status;
-- +goose StatementEnd
