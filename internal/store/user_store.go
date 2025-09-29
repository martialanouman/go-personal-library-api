package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	Id           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.hash = hash
	p.plaintext = &plaintext

	return nil
}

func (p *password) Matches(plaintext string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByToken(token, scope string) (*User, error)
}

type PostgresUserStore struct {
	db *pgxpool.Pool
}

func NewPostgresUserStore(db *pgxpool.Pool) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	ctx := context.Background()
	trx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer trx.Rollback(ctx)

	userInsertQuery := `
		INSERT INTO users (email, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err = s.db.QueryRow(
		ctx, userInsertQuery, user.Email, user.Name,
	).Scan(
		&user.Id, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	passwordInsertQuery := `
		INSERT INTO passwords (user_id, password_hash)
		VALUES ($1, $2)
	`

	_, err = s.db.Exec(ctx, passwordInsertQuery, user.Id, user.PasswordHash.hash)
	if err != nil {
		return err
	}

	err = trx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
		SELECT u.id, u.name, u.email, p.password_hash, u.created_at, u.updated_at
		FROM users u
		JOIN passwords p ON u.id = p.user_id
		WHERE email = $1
	`

	err := s.db.QueryRow(context.Background(), query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) GetUserByToken(token, scope string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
		SELECT u.id, u.name, u.email, u.created_at, u.updated_at
		FROM users u
		JOIN tokens t ON u.id = t.user_id
		WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > NOW()
	`

	hashedToken := sha256.Sum256([]byte(token))

	err := s.db.QueryRow(context.Background(), query, hashedToken[:], scope).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
