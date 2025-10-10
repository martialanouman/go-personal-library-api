package store

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martialanouman/personal-library/internal/utils"
)

const (
	ScopeAuth     = "auth"
	ScopeBooks    = "books"
	ScopeWishlist = "wishlist"
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
	RevokeAllTokens(userId, scope string) error
	GetTokenByHash(plaintext string) (*Token, error)
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

func (s *PostgresTokenStore) RevokeAllTokens(userId, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope LIKE '%' || $2 || '%'
	`

	_, err := s.db.Exec(context.Background(), query, userId, scope)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresTokenStore) GetTokenByHash(plaintext string) (*Token, error) {
	token := &Token{}
	hash := sha256.Sum256([]byte(plaintext))

	query := `
		SELECT user_id, scope, expiry, hash
		FROM tokens
		WHERE hash = $1
	`

	err := s.db.QueryRow(context.Background(), query, hash[:]).Scan(&token.UserId, &token.Scope, &token.Expiry, &token.Hash)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return token, nil
}
