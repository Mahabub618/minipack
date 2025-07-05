-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_items (
                             id UUID PRIMARY KEY,
                             order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                             product_id UUID NOT NULL,
                             product_name TEXT NOT NULL,
                             unit_price NUMERIC(10, 2) NOT NULL,
                             quantity INT NOT NULL CHECK (quantity > 0),
                             discount NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
                             total_price NUMERIC(10, 2) GENERATED ALWAYS AS (quantity * (unit_price - discount)) STORED,
                             created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
-- +goose StatementEnd
