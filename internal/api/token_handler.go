package api

import (
	"log"
	"net/http"

	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/middleware"
	"github.com/martialanouman/personal-library/internal/store"
)

type TokenHandler struct {
	store  store.TokenStore
	logger *log.Logger
}

func NewTokenHandler(store store.TokenStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		store:  store,
		logger: logger,
	}
}

func (h *TokenHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	err := h.store.RevokeAllTokens(user.Id, store.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: revoking token %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
