package handlers

import (
	"logger/data"
	"net/http"
)

var Model data.Model

func New(m data.Model) {
	Model = m
}

func WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload jsonPayload
	_ = readJSON(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := Model.LogEntry.Insert(event)

	if err != nil {
		errorJSON(w, err)
		return
	}
	response := jsonResponse{
		Error:   false,
		Message: "Log entry created",
		Data:    event,
	}

	writeJSON(w, http.StatusAccepted, response)
}
