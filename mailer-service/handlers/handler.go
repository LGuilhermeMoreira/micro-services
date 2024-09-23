package handlers

import (
	"encoding/json"
	"mailer-service/email"
	"net/http"
)

var mail email.Mail

func New(m email.Mail) {
	mail = m
}

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func SendMail(w http.ResponseWriter, r *http.Request) {
	var requestPayload mailMessage

	err := readJSON(w, r, &requestPayload)

	if err != nil {
		errorJSON(w, err)
		return
	}

	msg := email.Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = mail.SendSMTPMessage(msg)

	if err != nil {
		errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		errorJSON(w, err)
	}
}
