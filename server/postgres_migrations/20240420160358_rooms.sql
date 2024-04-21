-- +goose Up
-- +goose StatementBegin
CREATE TABLE rooms (
   id VARCHAR(255) PRIMARY KEY NOT NULL,
   name VARCHAR(255) NOT NULL UNIQUE,
   user_ids VARCHAR(255)[] DEFAULT '{}',
   created_by_id VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rooms;
-- +goose StatementEnd
