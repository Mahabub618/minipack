-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipped', 'cancelled', 'completed');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded');

CREATE TABLE orders (
                        id UUID PRIMARY KEY,
                        user_id UUID NOT NULL,
                        cart_id UUID,
                        order_number VARCHAR(20) UNIQUE NOT NULL,
                        status order_status NOT NULL DEFAULT 'pending',
                        payment_status payment_status NOT NULL DEFAULT 'pending',
                        payment_method TEXT,
                        shipping_address JSONB NOT NULL,
                        billing_address JSONB NOT NULL,
                        subtotal NUMERIC(10, 2) NOT NULL,
                        discount NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
                        tax NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
                        shipping_fee NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
                        total_amount NUMERIC(10, 2) NOT NULL,
                        currency VARCHAR(5) NOT NULL DEFAULT 'USD',
                        ordered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
