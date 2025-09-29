package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/store"
)

type UserMiddlewares struct {
	store store.UserStore
}

type UserContextString string

const (
	UserContextKey = UserContextString("user")
	AuthTokenType  = "Bearer"
)

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("could not get user from request context")
	}

	return user
}

func (m *UserMiddlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != AuthTokenType {
			helpers.WriteJson(w, http.StatusUnauthorized, helpers.Envelop{"error": "missing or invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := m.store.GetUserByToken(token, store.ScopeAuth)
		if err != nil {
			helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
			return
		}

		if user == nil {
			helpers.WriteJson(w, http.StatusUnauthorized, helpers.Envelop{"error": "missing or invalid authorization header"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (m *UserMiddlewares) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == store.AnonymousUser {
			helpers.WriteJson(w, http.StatusUnauthorized, helpers.Envelop{"error": "you must be logged in"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
