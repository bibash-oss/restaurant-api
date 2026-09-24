-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE menu_type AS ENUM ('BAR', 'KITCHEN');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE menu_items
ADD COLUMN IF NOT EXISTS menu_type menu_type NOT NULL DEFAULT 'KITCHEN';

CREATE INDEX IF NOT EXISTS idx_menu_items_menu_type ON menu_items(menu_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_menu_items_menu_type;
ALTER TABLE menu_items DROP COLUMN IF EXISTS menu_type;
DROP TYPE IF EXISTS menu_type;
-- +goose StatementEnd
