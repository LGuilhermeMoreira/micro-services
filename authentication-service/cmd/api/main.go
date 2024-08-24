package main

import (
	"authentication/config"
	database "authentication/connection"
	"authentication/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Auth service is on")
	conn := database.ConnectToDB()
	if conn == nil {
		log.Fatal("Could not connect to database")
	}
	app := config.NewConfig(conn, "3000")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", app.Webport),
		Handler: routes.GetMux(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
