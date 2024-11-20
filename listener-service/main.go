package main

import (
	"listener/event"
	"log"
)

func main() {
	// try to connect to rabbitmq
	conn, err := connet()

	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()
	log.Println("Connected to RabbitMQ")
	// start listening for messages
	log.Println("Listening for and consuming messages")
	// create consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		log.Panic(err)
	}
	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.Error"})
}
