package handlers

import (
	"authentication/config"
	"encoding/json"
	"errors"
	"net/http"
)

type Helper struct {
	cnfg config.Config
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != nil {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	return err
}

func errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	return writeJSON(w, statusCode, payload)
}
