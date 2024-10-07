package event

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

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
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}

}

func authEvent(payload Payload) error {
	return nil
}

func logEvent(payload Payload) error {
	return nil
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
