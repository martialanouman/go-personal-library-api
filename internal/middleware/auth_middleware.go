package middleware

import (
	"context"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/store"
)

type AuthMiddleware struct {
	UserStore  store.UserStore
	TokenStore store.TokenStore
	Logger     *log.Logger
}

type UserContextString string
type TokenContextKeyString string

const (
	UserContextKey  = UserContextString("user")
	TokenContextKey = TokenContextKeyString("token")
	AuthTokenType   = "Bearer"
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

func SetToken(r *http.Request, token *store.Token) *http.Request {
	ctx := context.WithValue(r.Context(), TokenContextKey, token)

	return r.WithContext(ctx)
}

func GetScope(r *http.Request) []string {
	token, ok := r.Context().Value(TokenContextKey).(*store.Token)
	if !ok {
		panic("could not get token from request context")
	}

	return strings.Split(token.Scope, ",")
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
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

		plaintextToken := headerParts[1]
		user, err := m.UserStore.GetUserByToken(plaintextToken, store.ScopeAuth)
		if err != nil {
			m.Logger.Printf("ERROR: getting user by token %v", err)
			helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
			return
		}

		if user == nil {
			helpers.WriteJson(w, http.StatusUnauthorized, helpers.Envelop{"error": "missing or invalid authorization header"})
			return
		}

		token, err := m.TokenStore.GetTokenByHash(plaintextToken)
		if err != nil {
			m.Logger.Printf("ERROR: getting token by hash %v", err)
			helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
			return
		}

		if token == nil {
			helpers.WriteJson(w, http.StatusUnauthorized, helpers.Envelop{"error": "missing or invalid authorization header"})
			return
		}

		r = SetUser(r, user)
		r = SetToken(r, token)
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user == store.AnonymousUser {
			helpers.WriteJson(w, http.StatusUnauthorized, helpers.Envelop{"error": "you must be logged in"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) RequireScope(next http.HandlerFunc, scope []string) http.HandlerFunc {
	return m.RequireUser(func(w http.ResponseWriter, r *http.Request) {
		tokenScope := GetScope(r)

		hasScopes := true
		for _, s := range scope {
			if !slices.Contains(tokenScope, s) {
				hasScopes = false
				break
			}
		}

		if !hasScopes {
			helpers.WriteJson(w, http.StatusForbidden, helpers.Envelop{"error": "you do not have the necessary permissions to access this resource"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
