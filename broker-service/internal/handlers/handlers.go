package handlers

import (
	"net/http"
)

func Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "entrou no broker",
	}

	_ = writeJSON(w, http.StatusOK, payload)
}
