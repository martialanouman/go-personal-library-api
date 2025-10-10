package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/middleware"
	"github.com/martialanouman/personal-library/internal/store"
)

type WishlistHandler struct {
	store  store.WishlistStore
	logger *log.Logger
}

func NewWishlistHandler(store store.WishlistStore, logger *log.Logger) WishlistHandler {
	return WishlistHandler{store: store, logger: logger}
}

type createWishRequest struct {
	Title     string  `json:"title"`
	Author    string  `json:"author"`
	Isbn      *string `json:"isbn"`
	BigBookID *int64  `json:"bb_id"`
	Priority  *string `json:"priority,omitempty"`
	Notes     *string `json:"notes,omitempty"`
}

func (req *createWishRequest) validate() map[string]string {
	errorMessages := make(map[string]string)

	if req.Title == "" {
		errorMessages["title"] = "title is required"
	}

	if req.Author == "" {
		errorMessages["author"] = "author is required"
	}

	if req.Isbn != nil && len(*req.Isbn) < 13 {
		errorMessages["isbn"] = "isbn must be 13 characters"
	}

	if req.Priority == nil {
		normal := "normal"
		req.Priority = &normal
	}

	priorities := []string{"low", "normal", "high"}
	if !slices.Contains(priorities, *req.Priority) {
		errorMessages["priority"] = "priority must be one of: low, normal or high"
	}

	return errorMessages
}

func (req *createWishRequest) toWish() *store.Wish {
	return &store.Wish{
		Title:     req.Title,
		Author:    &req.Author,
		Isbn:      req.Isbn,
		BigBookID: req.BigBookID,
		Priority:  *req.Priority,
		Acquired:  false,
		Notes:     req.Notes,
	}
}

func (h *WishlistHandler) HandleAddWish(w http.ResponseWriter, r *http.Request) {
	var req createWishRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: decoding payload %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	if validationMessages := req.validate(); len(validationMessages) > 0 {
		helpers.WriteJson(w, http.StatusUnprocessableEntity, helpers.Envelop{"errors": validationMessages})
		return
	}

	user := middleware.GetUser(r)
	wish := req.toWish()
	wish.UserID = user.ID
	if err := h.store.AddWish(wish); err != nil {
		h.logger.Printf("ERROR: adding wish %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelop{"wish": wish})
}

func (h *WishlistHandler) HandleDeleteWish(w http.ResponseWriter, r *http.Request) {

}

func (h *WishlistHandler) HandleMarkAsAcquired(w http.ResponseWriter, r *http.Request) {

}
