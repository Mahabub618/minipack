-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ALTER COLUMN order_number TYPE VARCHAR(30);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
ALTER COLUMN order_number TYPE VARCHAR(20);
-- +goose StatementEnd
