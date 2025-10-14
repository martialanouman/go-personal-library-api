package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Book struct {
	ID           string     `json:"id" db:"id"`
	UserId       string     `json:"user_id" db:"user_id"`
	Title        string     `json:"title" db:"title"`
	Author       string     `json:"author" db:"author"`
	Isbn         *string    `json:"isbn,omitempty" db:"isbn"`
	Description  *string    `json:"description,omitempty" db:"description"`
	CoverUrl     *string    `json:"cover_url,omitempty" db:"cover_url"`
	Genre        *string    `json:"genre,omitempty" db:"genre"`
	Status       string     `json:"status" db:"status"`
	Rating       byte       `json:"rating" db:"rating"`
	Notes        *string    `json:"notes,omitempty" db:"notes"`
	DateAdded    time.Time  `json:"date_added" db:"date_added"`
	DateStarted  *time.Time `json:"date_started,omitempty" db:"date_started"`
	DateFinished *time.Time `json:"date_finished,omitempty" db:"date_finished"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type BookStore interface {
	CreateBook(book *Book) error
	GetBooks(userId string, page, take int) ([]Book, error)
	GetBookById(id string) (*Book, error)
	UpdateBook(book *Book) error
	DeleteBook(id string) error
	GetBooksCount(userId string) (int, error)
}

type PostgresBookStore struct {
	db *pgxpool.Pool
}

func NewPostgresBookStore(db *pgxpool.Pool) *PostgresBookStore {
	return &PostgresBookStore{db}
}

func (s *PostgresBookStore) CreateBook(book *Book) error {
	query := `
		INSERT INTO books (user_id, title, author, isbn, description, cover_url, genre, status, rating, notes, date_added, date_started, date_finished)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(
		context.Background(), query,
		book.UserId,
		book.Title,
		book.Author,
		book.Isbn,
		book.Description,
		book.CoverUrl,
		book.Genre,
		book.Status,
		book.Rating,
		book.Notes,
		book.DateAdded,
		book.DateStarted,
		book.DateFinished,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresBookStore) GetBooks(userId string, page, take int) ([]Book, error) {
	query := "SELECT * FROM books WHERE user_id = $1 LIMIT $2 OFFSET $3 ORDER BY created_at DESC"

	rows, _ := s.db.Query(context.Background(), query, userId, take, (page-1)*take)
	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *PostgresBookStore) GetBookById(id string) (*Book, error) {
	var book *Book
	const query = "SELECT * FROM books WHERE id = $1"

	rows, _ := s.db.Query(context.Background(), query, id)
	book, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Book])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *PostgresBookStore) UpdateBook(book *Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, description = $4, cover_url = $5, genre = $6, status = $7, rating = $8, notes = $9, date_added = $10, date_started = $11, date_finished = $12, updated_at = NOW()
		WHERE id = $13
		RETURNING updated_at
	`

	err := s.db.QueryRow(
		context.Background(), query,
		book.Title,
		book.Author,
		book.Isbn,
		book.Description,
		book.CoverUrl,
		book.Genre,
		book.Status,
		book.Rating,
		book.Notes,
		book.DateAdded,
		book.DateStarted,
		book.DateFinished,
		book.ID,
	).Scan(&book.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresBookStore) DeleteBook(id string) error {
	query := "DELETE FROM books WHERE id = $1"
	commandTag, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *PostgresBookStore) GetBooksCount(userId string) (int, error) {
	var count int

	query := "SELECT COUNT(*) FROM books WHERE user_id = $1"
	err := s.db.QueryRow(context.Background(), query, userId).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
