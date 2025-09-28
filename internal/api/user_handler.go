package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/martialanouman/personal-library/internal/helpers"
	"github.com/martialanouman/personal-library/internal/store"
)

const (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type UserHandler struct {
	store  store.UserStore
	logger *log.Logger
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *registerUserRequest) validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}

	if rgx := regexp.MustCompile(emailRegex); !rgx.MatchString(r.Email) {
		return errors.New("email must be a valid email address")
	}

	if r.Name == "" {
		return errors.New("name is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func NewUserHandler(store store.UserStore, logger *log.Logger) UserHandler {
	return UserHandler{
		store,
		logger,
	}
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: decoding payload %v", err)
		helpers.WriteJson(w, http.StatusBadRequest, helpers.Envelop{"error": "invalid request payload"})
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Printf("ERROR: validating payload %v", err)
		helpers.WriteJson(w, http.StatusBadRequest, helpers.Envelop{"error": err.Error()})
		return
	}

	existingUser, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Printf("ERROR: checking existing user %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	if existingUser != nil {
		helpers.WriteJson(w, http.StatusConflict, helpers.Envelop{"error": "user with this email already exists"})
		return
	}

	user := &store.User{
		Email: req.Email,
		Name:  req.Name,
	}
	if err := user.PasswordHash.Set(req.Password); err != nil {
		h.logger.Printf("ERROR: hashing password %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	if err := h.store.CreateUser(user); err != nil {
		h.logger.Printf("ERROR: creating user %v", err)
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.Envelop{"error": "internal server error"})
		return
	}

	helpers.WriteJson(w, http.StatusCreated, helpers.Envelop{"user": user})
}
