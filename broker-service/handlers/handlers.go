package handlers

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/rabbitmq/amqp091-go"
)

type Controller struct {
	conn *amqp091.Connection
}

func NewController(conn *amqp091.Connection) *Controller {
	return &Controller{
		conn: conn,
	}
}

func (c *Controller) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "entrou no broker",
	}

	_ = writeJSON(w, http.StatusOK, payload)
}

func (c *Controller) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := readJSON(w, r, &requestPayload)
	if err != nil {
		errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		c.authenticate(w, requestPayload.Auth)
	case "log":
		// c.logger(w, requestPayload.Log)
		c.loggerRabbitMQ(w, requestPayload.Log)
	case "mail":
		c.mail(w, requestPayload.Mail)
	default:
		errorJSON(w, errors.New("unknown action"))
	}
}

func (c *Controller) authenticate(w http.ResponseWriter, a AuthPayload) {
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

func (c *Controller) logger(w http.ResponseWriter, entry LogPayload) {
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
func (c *Controller) mail(w http.ResponseWriter, payload MailPayload) {
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
	json.NewEncoder(w).Encode(responsePayload)
}

func (c *Controller) loggerRabbitMQ(w http.ResponseWriter, payload LogPayload) {
	err := c.pushToQueue(payload)
	if err != nil {
		errorJSON(w, err)
		return
	}
	jResp := jsonResponse{
		Error:   false,
		Message: "logged via RabbitMQ",
		Data:    nil,
	}

	writeJSON(w, http.StatusAccepted, jResp)
}

func (c *Controller) pushToQueue(payload LogPayload) error {
	emitter, err := event.NewEmitter(c.conn)
	if err != nil {
		return err
	}
	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}
	return emitter.Push(string(j), "log.INFO")
}
