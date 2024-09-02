package main

import (
	"context"
	"fmt"
	"log"
	"logger/config"
	"logger/routes"
	"net/http"
	"time"
)

func main() {
	cnfg, err := config.NewConfig("80",
		"5001",
		"50001",
		"mongodb://mongodb:27017")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*150)
	defer cancel()
	defer func() {
		if err = cnfg.MongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cnfg.WebPort),
		Handler: routes.GetMux(*cnfg),
	}
	log.Println("logger service started")
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
