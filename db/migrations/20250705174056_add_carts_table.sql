-- +goose Up
-- +goose StatementBegin
CREATE TYPE cart_status AS ENUM ('active', 'abandoned', 'converted');
CREATE TABLE carts (
                       id UUID PRIMARY KEY,
                       user_id UUID NOT NULL,
                       status cart_status NOT NULL DEFAULT 'active',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       expires_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE carts;
-- +goose StatementEnd
