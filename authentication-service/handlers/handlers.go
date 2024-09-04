package handlers

import (
	"authentication/data"
	"bytes"
	"database/sql"
	"encoding/json"
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

	err = logRegister("auth", fmt.Sprintf("%v is logged", user.Email))

	if err != nil {
		errorJSON(w, err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	writeJSON(w, http.StatusAccepted, payload)
}

func logRegister(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")

	if err != nil {
		return err
	}

	url := "http://logger-service/log"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return errors.New("request not accepted")
	}

	return nil
}
