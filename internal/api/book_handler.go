package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/middleware"
	"github.com/martialanouman/personal-library/internal/store"
)

type BookHandler struct {
	store  store.BookStore
	logger *log.Logger
}

type createBookRequest struct {
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Isbn         *string `json:"isbn,omitempty"`
	Description  *string `json:"description,omitempty"`
	CoverUrl     *string `json:"cover_url,omitempty"`
	Genre        *string `json:"genre,omitempty"`
	Status       string  `json:"status"`
	Rating       byte    `json:"rating"`
	Notes        *string `json:"notes,omitempty"`
	DateStarted  *string `json:"date_started,omitempty"`
	DateFinished *string `json:"date_finished,omitempty"`
	DateAdded    *string `json:"date_added,omitempty"`
}

func (r *createBookRequest) validate() map[string]string {
	errorMessages := make(map[string]string)

	if r.Title == "" {
		errorMessages["title"] = "title is required"
	}

	if r.Author == "" {
		errorMessages["author"] = "author is required"
	}

	if r.Status == "" {
		errorMessages["status"] = "status is required"
	} else {
		validStatuses := []string{"to_read", "reading", "read"}
		isValidStatus := slices.Contains(validStatuses, r.Status)

		if !isValidStatus {
			errorMessages["status"] = "status must be one of: to_read, reading, read"
		}
	}

	if r.Rating < 1 || r.Rating > 5 {
		errorMessages["rating"] = "rating must be between 1 and 5"
	}

	if r.DateAdded != nil {
		_, err := time.Parse(time.DateOnly, *r.DateAdded)
		if err != nil {
			errorMessages["date_added"] = "date_added must be in YYYY-MM-DD format"
		}
	}

	if r.DateStarted != nil {
		_, err := time.Parse(time.DateOnly, *r.DateStarted)
		if err != nil {
			errorMessages["date_started"] = "date_started must be in YYYY-MM-DD format"
		}
	}

	if r.DateFinished != nil {
		_, err := time.Parse(time.DateOnly, *r.DateFinished)
		if err != nil {
			errorMessages["date_finished"] = "date_finished must be in YYYY-MM-DD format"
		}
	}

	return errorMessages
}

func (r *createBookRequest) toBook() *store.Book {
	var (
		dateAdded    time.Time = time.Now()
		dateFinished *time.Time
		DateStarted  *time.Time
	)
	if r.DateAdded != nil {
		parsedDate, _ := time.Parse(time.DateOnly, *r.DateAdded)
		dateAdded = parsedDate
	}

	if r.DateFinished != nil {
		parsedDate, _ := time.Parse(time.DateOnly, *r.DateFinished)
		dateFinished = &parsedDate
	}

	if r.DateStarted != nil {
		parsedDate, _ := time.Parse(time.DateOnly, *r.DateStarted)
		DateStarted = &parsedDate
	}

	return &store.Book{
		Title:        r.Title,
		Author:       r.Author,
		Isbn:         r.Isbn,
		Description:  r.Description,
		CoverUrl:     r.CoverUrl,
		Genre:        r.Genre,
		Status:       r.Status,
		Rating:       r.Rating,
		Notes:        r.Notes,
		DateStarted:  DateStarted,
		DateFinished: dateFinished,
		DateAdded:    dateAdded,
	}
}

func NewBookHandler(store store.BookStore, logger *log.Logger) BookHandler {
	return BookHandler{store: store, logger: logger}
}

func (h *BookHandler) HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	books, err := h.store.GetBooks(user.Id)
	if err != nil {
		h.logger.Printf("ERROR: getting books %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"books": books})

}

func (h *BookHandler) HandlerCreateBook(w http.ResponseWriter, r *http.Request) {
	var req createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: decoding create book request %v", err)
		helpers.WriteJson(w, http.StatusBadRequest, helpers.Envelop{"error": "invalid request payload"})
		return
	}

	if validationErrors := req.validate(); len(validationErrors) > 0 {
		helpers.WriteJson(w, http.StatusUnprocessableEntity, helpers.Envelop{"errors": validationErrors})
		return
	}

	user := middleware.GetUser(r)
	book := req.toBook()
	book.UserId = user.Id

	if err := h.store.CreateBook(book); err != nil {
		h.logger.Printf("ERROR: creating book %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelop{"book": book})
}

func (h *BookHandler) HandleGetBookById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.Envelop{"error": "invalid book id"})
		return
	}

	book, err := h.store.GetBookById(id)
	if err != nil {
		h.logger.Printf("ERROR: getting book by id %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	if book == nil {
		helpers.WriteJson(w, http.StatusNotFound, helpers.Envelop{"error": "book not found"})
		return
	}

	helpers.WriteJson(w, http.StatusOK, helpers.Envelop{"book": book})
}
