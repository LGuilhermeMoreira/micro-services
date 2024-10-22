package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "entrou no broker",
	}

	_ = writeJSON(w, http.StatusOK, payload)
}

func HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		authenticate(w, requestPayload.Auth)
	case "log":
		logger(w, requestPayload.Log)
	case "mail":
		mail(w, requestPayload.Mail)
	default:
		errorJSON(w, errors.New("unknown action"))
	}
}

func authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json to send to auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		errorJSON(w, err)
		return
	}

	client := http.DefaultClient

	response, err := client.Do(request)
	if err != nil {
		errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Println("wrong status of", response.StatusCode)
		errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variabel read response.Body into
	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		errorJSON(w, err)
	}

	if jsonFromService.Error {
		errorJSON(w, err, http.StatusUnauthorized)
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	writeJSON(w, http.StatusAccepted, payload)

}

func logger(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		errorJSON(w, err, http.StatusBadRequest)
		return
	}

	loggerServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", loggerServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		errorJSON(w, err, http.StatusBadGateway)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		errorJSON(w, err, http.StatusBadGateway)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		errorJSON(w, errors.New("status not accepted"), http.StatusNotAcceptable)
		return
	}

	var respData map[string]any

	if err = json.NewDecoder(response.Body).Decode(&respData); err != nil {
		errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	msg, _ := respData["message"].(string)

	var resp jsonResponse
	resp.Error = false
	resp.Message = msg
	resp.Data = respData["data"]

	writeJSON(w, http.StatusAccepted, resp)
}

func mail(w http.ResponseWriter, payload MailPayload) {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	mailServiceUrl := "http://mailer-service/send"
	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		errorJSON(w, err, http.StatusBadGateway)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println(err)
		errorJSON(w, err, http.StatusBadGateway)
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		log.Println(response.StatusCode)
		errorJSON(w, errors.New("status not accepted in mail service"), http.StatusNotAcceptable)
		return
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Mail sent!"
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(responsePayload)
}
