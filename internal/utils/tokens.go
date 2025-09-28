package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func GenerateToken(ttl time.Duration) (*Token, error) {
	token := &Token{
		Expiry: time.Now().Add(ttl),
	}

	emptyBytes := make([]byte, 32)
	if _, err := rand.Read(emptyBytes); err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
