package main

import (
	"fmt"
	"log"
	"mailer-service/config"
	"mailer-service/routes"
	"net/http"
)

func main() {
	log.Println("starting mail service")

	cnfg := config.NewConfig("80")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cnfg.WebPort),
		Handler: routes.GetMux(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
