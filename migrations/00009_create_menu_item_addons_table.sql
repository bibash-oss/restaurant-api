-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS menu_item_addons (
    menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    addon_id UUID NOT NULL REFERENCES addons(id) ON DELETE CASCADE,
    PRIMARY KEY (menu_item_id, addon_id)
);

CREATE INDEX IF NOT EXISTS idx_menu_item_addons_item ON menu_item_addons(menu_item_id);
CREATE INDEX IF NOT EXISTS idx_menu_item_addons_addon ON menu_item_addons(addon_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menu_item_addons;
-- +goose StatementEnd
