package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	for _, s := range topics {
		if err := ch.QueueBind(q.Name, s, "logs_topic", false, nil); err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	for d := range messages {
		var payload Payload
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("Invalid message format: %s", err)
			continue
		}
		go handlePayload(payload)
	}

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// authenticate

	// you can have as many cases as you want, as long as you write the logic

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request to log service: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("log service returned status %d: %s", response.StatusCode, string(body))
	}

	return nil
}
