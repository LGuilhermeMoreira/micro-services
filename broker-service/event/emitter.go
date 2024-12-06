package event

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return declareExchange(ch)
}

func (e *Emitter) Push(event string, severity string) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	log.Println("Push to channel")
	err = ch.Publish(
		"log_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func NewEmitter(conn *amqp.Connection) (*Emitter, error) {
	emmiter := Emitter{
		connection: conn,
	}
	err := emmiter.setup()
	if err != nil {
		return nil, err
	}
	return &emmiter, nil
}
