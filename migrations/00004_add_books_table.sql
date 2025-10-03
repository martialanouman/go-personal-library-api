-- +goose Up
-- +goose StatementBegin
CREATE TYPE BOOK_STATUS AS ENUM ('to_read', 'reading', 'read');
CREATE TABLE IF NOT EXISTS books (
    id UUID DEFAULT UUIDV7() PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(13) UNIQUE,
    description TEXT,
    cover_url VARCHAR(255),
    genre VARCHAR(255),
    status BOOK_STATUS DEFAULT 'to_read',
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    notes TEXT,
    date_added TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    date_started TIMESTAMP WITH TIME ZONE,
    date_finished TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS books;
DROP TYPE IF EXISTS BOOK_STATUS;
-- +goose StatementEnd
