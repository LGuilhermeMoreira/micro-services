package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	// try to connect to rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()
	log.Println("Connected to RabbitMQ")
	// start listening for messages

	// create consumer

	// watch the queue and consume events
}
