-- +goose Up
-- +goose StatementBegin
-- Add table_name to orders
ALTER TABLE orders
ADD COLUMN IF NOT EXISTS table_name VARCHAR(255);

-- If table_id exists in orders, copy over table number if tables table exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'orders' AND column_name = 'table_id'
    ) AND EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_name = 'tables'
    ) THEN
        UPDATE orders o
        SET table_name = t.number::text
        FROM tables t
        WHERE o.table_id = t.id AND (o.table_name IS NULL OR o.table_name = '');
    END IF;
END $$;

-- Drop table_id foreign key constraint and column from orders
ALTER TABLE orders
DROP CONSTRAINT IF EXISTS orders_table_id_fkey;

DROP INDEX IF EXISTS idx_orders_table_id;

ALTER TABLE orders
DROP COLUMN IF EXISTS table_id;

-- Add table_name to payment_sessions
ALTER TABLE payment_sessions
ADD COLUMN IF NOT EXISTS table_name VARCHAR(255);

-- If table_id exists in payment_sessions, copy over
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'payment_sessions' AND column_name = 'table_id'
    ) AND EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_name = 'tables'
    ) THEN
        UPDATE payment_sessions p
        SET table_name = t.number::text
        FROM tables t
        WHERE p.table_id = t.id AND (p.table_name IS NULL OR p.table_name = '');
    END IF;
END $$;

-- Drop table_id foreign key constraint and column from payment_sessions
ALTER TABLE payment_sessions
DROP CONSTRAINT IF EXISTS payment_sessions_table_id_fkey;

DROP INDEX IF EXISTS idx_payment_sessions_table_id;

ALTER TABLE payment_sessions
DROP COLUMN IF EXISTS table_id;

-- Drop tables table
DROP TABLE IF EXISTS tables;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    number INT NOT NULL,
    capacity INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    CONSTRAINT uq_restaurant_table_number UNIQUE (restaurant_id, number)
);

ALTER TABLE orders
ADD COLUMN IF NOT EXISTS table_id UUID REFERENCES tables(id) ON DELETE RESTRICT;

ALTER TABLE orders
DROP COLUMN IF EXISTS table_name;

ALTER TABLE payment_sessions
ADD COLUMN IF NOT EXISTS table_id UUID REFERENCES tables(id) ON DELETE RESTRICT;

ALTER TABLE payment_sessions
DROP COLUMN IF EXISTS table_name;
-- +goose StatementEnd
