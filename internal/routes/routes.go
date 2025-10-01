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
		r.Post("/login", app.UserHandler.HandleLogin)

		r.Group(func(r chi.Router) {
			r.Use(app.AuthMiddleware.Authenticate)

			r.Get("/me", app.AuthMiddleware.RequireUser(app.UserHandler.HandleMe))
			r.Put("/password", app.AuthMiddleware.RequireUser(app.UserHandler.HandleUpdatePassword))
			r.Delete("/logout", app.AuthMiddleware.RequireUser(app.TokenHandler.HandleLogout))
		})
	})

	return r
}
