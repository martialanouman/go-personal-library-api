package app

import (
	"database/sql"
	"log"
	"os"
)

type Application struct {
	db     *sql.DB
	Logger *log.Logger
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "[APP] ", log.Ldate|log.Ltime)

	return &Application{
		Logger: logger,
		db:     nil,
	}, nil
}
