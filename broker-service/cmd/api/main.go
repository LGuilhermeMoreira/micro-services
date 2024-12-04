package main

import (
	"broker/config"
	"broker/routes"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cnfg := config.NewConfig("80")
	conn, err := connect()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to rabbitMQ")
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cnfg.WebPort),
		Handler: routes.GetMux(conn),
	}
	log.Println("Broker service is on")

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
