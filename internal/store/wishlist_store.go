package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Wish struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Author    *string   `json:"author,omitempty" db:"author"`
	Isbn      *string   `json:"isbn,omitempty" db:"isbn"`
	BigBookID *int64    `json:"bb_id,omitempty" db:"big_book_id"`
	Priority  string    `json:"priority" db:"priority"`
	Acquired  bool      `json:"acquired" db:"acquired"`
	Notes     *string   `json:"notes" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WishlistStore interface {
	AddWish(wish *Wish) error
	GetWishById(id string) (*Wish, error)
	GetWishes(userId string, page, take int) ([]Wish, error)
	DeleteWishById(id string) error
	MarkAsAcquired(id string) error
	GetWishesCount(userId string) (int, error)
}

type PostgresWishlistStore struct {
	db *pgxpool.Pool
}

func NewPostgresWishlistStore(db *pgxpool.Pool) *PostgresWishlistStore {
	return &PostgresWishlistStore{db}
}

func (s *PostgresWishlistStore) AddWish(wish *Wish) error {

	query := `
		INSERT INTO wishlists (user_id, title, author, isbn, big_book_id, priority, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := s.db.QueryRow(context.Background(),
		query,
		wish.UserID,
		wish.Title,
		wish.Author,
		wish.Isbn,
		wish.BigBookID,
		wish.Priority,
		wish.Notes,
	).Scan(&wish.ID, &wish.CreatedAt)
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

func (s *PostgresWishlistStore) DeleteWishById(id string) error {
	query := `DELETE FROM wishlists WHERE id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresWishlistStore) GetWishes(userId string, page, take int) ([]Wish, error) {
	query := `
		SELECT * 
		FROM wishlists 
		WHERE user_id = $1 AND acquired = FALSE
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`
	offset := (page - 1) * take

	rows, err := s.db.Query(context.Background(), query, userId, take, offset)
	if err != nil {
		return nil, err
	}

	wishes, err := pgx.CollectRows(rows, pgx.RowToStructByName[Wish])
	if err != nil {
		return nil, err
	}

	return wishes, nil
}

func (s *PostgresWishlistStore) GetWishesCount(userId string) (int, error) {
	query := "SELECT COUNT(*) FROM wishlists WHERE user_id = $1 AND acquired = FALSE"

	var count int
	err := s.db.QueryRow(context.Background(), query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *PostgresWishlistStore) MarkAsAcquired(id string) error {
	ctx := context.Background()

	trx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer trx.Rollback(ctx)

	query := `
		UPDATE wishlists 
		SET acquired = TRUE, updated_at = NOW()
		WHERE id = $1
		RETURNING *
	`

	rows, err := trx.Query(ctx, query, id)
	if err != nil {
		return err
	}

	wish, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Wish])
	if err != nil {
		return err
	}

	if !wish.Acquired {
		return nil
	}

	insertBookQuery := `
		INSERT INTO books (user_id, title, author, isbn, notes, rating)
		VALUES ($1, $2, $3, $4, $5, 1)
	`

	_, err = trx.Exec(ctx, insertBookQuery, wish.UserID, wish.Title, wish.Author, wish.Isbn, wish.Notes)
	if err != nil {
		return err
	}

	err = trx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
