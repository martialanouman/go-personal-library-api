package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/martialanouman/personal-library/internal/helpers"
)

type Application struct {
	db     *sql.DB
	Logger *log.Logger
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "[APP]: ", log.Ldate|log.Ltime)

	return &Application{
		Logger: logger,
		db:     nil,
	}, nil
}

func (a *Application) Health(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"status": "up"})
}
