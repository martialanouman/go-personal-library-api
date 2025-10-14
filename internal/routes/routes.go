package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/martialanouman/personal-library/internal/app"
	"github.com/martialanouman/personal-library/internal/store"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", app.Health)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.UserHandler.HandleRegisterUser)
			r.Post("/login", app.UserHandler.HandleLogin)

			r.Group(func(r chi.Router) {
				r.Use(app.AuthMiddleware.Authenticate)

				r.Get("/me", app.AuthMiddleware.RequireScope(app.UserHandler.HandleMe, []string{store.ScopeAuth}))
				r.Put("/password", app.AuthMiddleware.RequireScope(app.UserHandler.HandleUpdatePassword, []string{store.ScopeAuth}))
				r.Delete("/logout", app.AuthMiddleware.RequireUser(app.TokenHandler.HandleLogout))
			})
		})

		r.Route("/books", func(r chi.Router) {
			r.Use(app.AuthMiddleware.Authenticate)

			r.With(app.UtilsMiddleware.GetPagination).Get("/", app.AuthMiddleware.RequireScope(app.BookHandler.HandleGetBooks, []string{store.ScopeBooks}))
			r.Post("/", app.AuthMiddleware.RequireScope(app.BookHandler.HandlerCreateBook, []string{store.ScopeBooks}))
			r.Post("/import/{bbId}", app.AuthMiddleware.RequireScope(app.BookHandler.HandleAddBookByISBN, []string{store.ScopeBooks}))
			r.Get("/{id}", app.AuthMiddleware.RequireScope(app.BookHandler.HandleGetBookById, []string{store.ScopeBooks}))
			r.Put("/{id}", app.AuthMiddleware.RequireScope(app.BookHandler.HandleUpdateBook, []string{store.ScopeBooks}))
			r.Delete("/{id}", app.AuthMiddleware.RequireScope(app.BookHandler.HandleDeleteBook, []string{store.ScopeBooks}))
		})

		r.Route("/wishes", func(r chi.Router) {
			r.Use(app.AuthMiddleware.Authenticate)

			r.Post("/", app.AuthMiddleware.RequireScope(app.WishlistHandler.HandleAddWish, []string{"wishlist"}))
			r.Delete("/{id}", app.AuthMiddleware.RequireScope(app.WishlistHandler.HandleDeleteWish, []string{"wishlist"}))
			r.Put("/{id}/acquire", app.AuthMiddleware.RequireScope(app.WishlistHandler.HandleMarkAsAcquired, []string{"wishlist"}))
			r.With(app.UtilsMiddleware.GetPagination).Get("/", app.AuthMiddleware.RequireScope(app.WishlistHandler.HandleGetWishes, []string{store.ScopeWishlist}))
		})
	})

	return r
}
