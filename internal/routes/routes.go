package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/martialanouman/personal-library/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.Health)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.UserHandler.HandleRegisterUser)

	})

	return r
}
