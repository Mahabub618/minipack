-- +goose Up
-- +goose StatementBegin
CREATE TABLE anonymous_users (
    id UUID PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE anonymous_users;
-- +goose StatementEnd
