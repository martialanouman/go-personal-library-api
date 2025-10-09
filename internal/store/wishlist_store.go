package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Wish struct {
	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    *string   `json:"author,omitempty" db:"author"`
	Isbn      *string   `json:"isbn,omitempty" db:"isbn"`
	BigBookID *int64    `json:"bb_id,omitempty" db:"big_book_id"`
	Priority  string    `json:"priority" db:"priority"`
	Acquired  bool      `json:"acquired" db:"acquired"`
	Notes     *string   `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type WishlistStore interface {
	AddWish(wish *Wish) error
	GetWishById(id string) (*Wish, error)
	DeleteWish(id string) error
	MarkAsAcquired(id string) error
}

type PostgresWishlistStore struct {
	db *pgxpool.Pool
}

func NewPostgresWishlistStore(db *pgxpool.Pool) *PostgresBookStore {
	return &PostgresBookStore{db}
}

func (s *PostgresWishlistStore) AddWish(wish *Wish) error {

	query := `
		INSERT INTO wishlists (title, author, isbn, big_book_id, priority, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(context.Background(),
		query,
		wish.Title,
		wish.Author,
		wish.Isbn,
		wish.BigBookID,
		wish.Priority,
		wish.Notes,
	).Scan(&wish.ID, wish.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresWishlistStore) GetWishById(id string) (*Wish, error) {
	query := `SELECT * FROM wishlists WHERE id = $1`

	rows, err := s.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}

	wish, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Wish])
	if err != nil {
		return nil, err
	}

	return wish, nil
}

func (s *PostgresWishlistStore) DeleteWish(id string) error {
	query := `DELETE FROM wishlists WHERE id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresWishlistStore) MarkAsAcquired(id string) error {
	query := `
		UPDATE wishlists 
		SET (acquired = TRUE, updated_at = NOW()) 
		WHERE id = $1
	`

	_, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	return nil
}
