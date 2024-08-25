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
