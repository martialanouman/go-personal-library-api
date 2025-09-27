-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS passwords (
    id UUID DEFAULT UUIDV7() PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS passwords;
-- +goose StatementEnd
