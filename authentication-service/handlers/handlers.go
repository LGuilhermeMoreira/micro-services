package handlers

import (
	"authentication/data"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

type handlerConfig struct {
	DB *sql.DB
}

func NewhandlerConfig(db *sql.DB) *handlerConfig {
	return &handlerConfig{
		DB: db,
	}
}

func (h handlerConfig) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &requestPayload)

	if err != nil {
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	model := data.New(h.DB)

	user, err := model.User.GetByEmail(requestPayload.Email)
	if err != nil {
		errorJSON(w, errors.New("invalid credentials"), http.StatusNotFound)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	writeJSON(w, http.StatusAccepted, payload)
}
