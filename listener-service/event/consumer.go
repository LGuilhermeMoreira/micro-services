package event

import (
	"bytes"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()

	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	queue, err := declareRandomQueue(channel)
	if err != nil {
		return err
	}
	for _, topic := range topics {
		err = channel.QueueBind(queue.Name, topic, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}
	messages, err := channel.Consume(queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for message := range messages {
			var payload Payload
			_ = json.Unmarshal(message.Body, &payload)
			go handlePayload(payload)
		}
	}()
	log.Printf(" [*] Waiting for messages")
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//log
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		err := authEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "mail":
		err := mailEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}

}

func mailEvent(payload Payload) error {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")
	mailServiceUrl := "http://mailer-service/send"
	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		log.Println(response.StatusCode)
		return errors.New("status not accepted in mail service")
	}

	return nil
}

func authEvent(payload Payload) error {
	// create some json to send to auth microservice
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := http.DefaultClient

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("authentication is unauthorized")
	} else if response.StatusCode != http.StatusAccepted {
		log.Println("wrong status of", response.StatusCode)
		return errors.New("authentication failed")
	}
	return nil
}

func logEvent(payload Payload) error {
	jsonData, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	loggerServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", loggerServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
