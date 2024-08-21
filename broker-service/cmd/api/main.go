package main

import (
	"broker/config"
	"broker/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {

	cnfg := config.NewConfig("80")

	log.Println("Broker service is on")

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cnfg.WebPort),
		Handler: routes.GetMux(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
