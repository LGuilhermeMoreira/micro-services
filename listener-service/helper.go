package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func connet() (*amqp.Connection, error) {
	var counts int64
	var backOff = 2 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ is not ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}
		fmt.Println("Waiting for connection...")
		time.Sleep(backOff)
	}

	return connection, nil
}
