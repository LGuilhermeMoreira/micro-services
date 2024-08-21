package main

import (
	"authentication/config"
	"authentication/data"
	"authentication/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Auth service is on")
	//todo connect with database
	app := config.NewConfig(nil, "3000", data.Models{})
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", app.Webport),
		Handler: routes.GetMux(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
