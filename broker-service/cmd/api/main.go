package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct{}

var webPort = "3000"

func main() {
	app := Config{}

	log.Printf("server on %s 🔥\n", webPort)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
