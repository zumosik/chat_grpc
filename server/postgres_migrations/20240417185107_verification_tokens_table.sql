-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_confirm_tokens (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_confirm_tokens;
-- +goose StatementEnd
