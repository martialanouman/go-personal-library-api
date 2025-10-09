-- +goose Up
-- +goose StatementBegin
CREATE TYPE WISH_PRIORITY AS ENUM ('low', 'normal', 'high');
CREATE TABLE IF NOT EXISTS wishlists (
    id UUID DEFAULT UUIDV7() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    isbn VARCHAR(13),
    big_book_id BIGINT,
    priority WISH_PRIORITY DEFAULT 'low',
    acquired BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

);
COMMENT ON COLUMN wishlists.big_book_id IS 'ID from the Big Book API, if available';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wishlists;
DROP TYPE IF EXISTS WISH_PRIORITY;
-- +goose StatementEnd
    