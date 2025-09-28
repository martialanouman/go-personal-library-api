package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martialanouman/personal-library/internal/utils"
)

const (
	ScopeAuth = "auth"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	UserId    string    `json:"-"`
	Scope     string    `json:"-"`
}

type TokenStore interface {
	CreateToken(token *Token) error
}

type PostgresTokenStore struct {
	db *pgxpool.Pool
}

func NewPostgresTokenStore(db *pgxpool.Pool) *PostgresTokenStore {
	return &PostgresTokenStore{
		db,
	}
}

func (s *PostgresTokenStore) CreateToken(token *Token) error {
	genToken, err := utils.GenerateToken(24 * time.Hour)
	if err != nil {
		return err
	}

	token.Hash = genToken.Hash
	token.Plaintext = genToken.Plaintext
	token.Expiry = genToken.Expiry
	token.Scope = ScopeAuth

	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	_, err = s.db.Exec(context.Background(), query, token.Hash, token.UserId, token.Expiry, token.Scope)
	if err != nil {
		return err
	}

	return nil
}
