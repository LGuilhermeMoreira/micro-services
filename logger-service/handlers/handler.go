package handlers

import (
	"logger/data"
	"net/http"
)

type Logger struct {
	Model data.Model
}

func NewLogger(m data.Model) *Logger {
	return &Logger{
		Model: m,
	}
}

func (l *Logger) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload jsonPayload
	_ = readJSON(w, r, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := l.Model.LogEntry.Insert(event)

	if err != nil {
		errorJSON(w, err)
		return
	}
	response := jsonResponse{
		Error:   false,
		Message: "Log entry created",
		Data:    event,
	}

	writeJSON(w, http.StatusCreated, response)
}
