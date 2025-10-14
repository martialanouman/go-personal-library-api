package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
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
	wishID := chi.URLParam(r, "id")
	if wishID == "" {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.Envelop{"error": "invalid wish id"})
		return
	}

	user := middleware.GetUser(r)
	wish, err := h.store.GetWishById(wishID)
	if err != nil {
		h.logger.Printf("ERROR: getting wish %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	if user.ID != wish.UserID {
		helpers.WriteJson(w, http.StatusForbidden, helpers.Envelop{"error": "you are not allowed to perform this action on this resource"})
		return
	}

	if err := h.store.DeleteWishById(wish.ID); err != nil {
		h.logger.Printf("ERROR: deleting wish %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WishlistHandler) HandleMarkAsAcquired(w http.ResponseWriter, r *http.Request) {
	wishID := chi.URLParam(r, "id")
	if wishID == "" {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.Envelop{"error": "invalid wish id"})
		return
	}

	user := middleware.GetUser(r)
	wish, err := h.store.GetWishById(wishID)
	if err != nil {
		h.logger.Printf("ERROR: getting wish %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	if user.ID != wish.UserID {
		helpers.WriteJson(w, http.StatusForbidden, helpers.Envelop{"error": "you are not allowed to perform this action on this resource"})
		return
	}

	if err := h.store.MarkAsAcquired(wishID); err != nil {
		h.logger.Printf("ERROR: mark as acquired wish %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WishlistHandler) HandleGetWishes(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	pagination := middleware.GetPagination(r)

	wishes, err := h.store.GetWishes(user.ID, pagination.Page, pagination.Take)
	if err != nil {
		h.logger.Printf("ERROR: getting wishes %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	count, err := h.store.GetWishesCount(user.ID)
	if err != nil {
		h.logger.Printf("ERROR: getting wishes count %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	helpers.WriteJson(
		w, http.StatusOK,
		helpers.Envelop{"wishes": wishes, "page": pagination.Page, "take": pagination.Take, "count": count},
	)
}
