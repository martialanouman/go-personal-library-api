package store

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Open() (*pgxpool.Pool, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	conn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("ERROR: db open %v", err)
	}

	return conn, nil
}
