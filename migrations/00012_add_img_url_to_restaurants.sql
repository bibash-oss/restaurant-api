-- +goose Up
-- +goose StatementBegin
ALTER TABLE restaurants
ADD COLUMN IF NOT EXISTS img_url VARCHAR(500) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE restaurants
DROP COLUMN IF EXISTS img_url;
-- +goose StatementEnd
