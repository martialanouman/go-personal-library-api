package app

import (
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martialanouman/personal-library/internal/api"
	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/middleware"
	"github.com/martialanouman/personal-library/internal/store"
)

type Application struct {
	Db              *pgxpool.Pool
	Logger          *log.Logger
	AuthMiddleware  middleware.AuthMiddleware
	UtilsMiddleware middleware.UtilsMiddleware
	UserHandler     api.UserHandler
	TokenHandler    api.TokenHandler
	BookHandler     api.BookHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "[APP]: ", log.Ldate|log.Ltime)

	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	userStore := store.NewPostgresUserStore(db)
	tokenStore := store.NewPostgresTokenStore(db)
	bookStore := store.NewPostgresBookStore(db)

	return &Application{
		Logger:          logger,
		Db:              db,
		AuthMiddleware:  middleware.NewAuthMiddleware(userStore, tokenStore, logger),
		UtilsMiddleware: middleware.NewUtilsMiddleware(),
		UserHandler:     api.NewUserHandler(userStore, tokenStore, logger),
		TokenHandler:    api.NewTokenHandler(tokenStore, logger),
		BookHandler:     api.NewBookHandler(bookStore, logger),
	}, nil
}

func (a *Application) Health(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"status": "up"})
}
