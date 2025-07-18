-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items
DROP COLUMN product_id;

ALTER TABLE cart_items
    ADD COLUMN product_id INTEGER NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cart_items
DROP COLUMN product_id;

ALTER TABLE cart_items
    ADD COLUMN product_id UUID NOT NULL;
-- +goose StatementEnd
