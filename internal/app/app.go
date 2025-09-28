package app

import (
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martialanouman/personal-library/internal/api"
	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/store"
)

type Application struct {
	Db          *pgxpool.Pool
	Logger      *log.Logger
	UserHandler api.UserHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "[APP]: ", log.Ldate|log.Ltime)

	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	return &Application{
		Logger:      logger,
		Db:          db,
		UserHandler: api.NewUserHandler(store.NewPostgresUserStore(db), logger),
	}, nil
}

func (a *Application) Health(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"status": "up"})
}
