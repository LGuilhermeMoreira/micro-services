package main

import (
	"broker/config"
	"broker/internal/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {

	cnfg := config.NewConfig("3000")

	log.Printf("server on %s 🔥\n", cnfg.WebPort)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cnfg.WebPort),
		Handler: routes.GetMux(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
