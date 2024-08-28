package main

import (
	"context"
	"logger/cmd/config"
	"time"
)

func main() {
	cnfg, err := config.NewConfig("80",
		"5001",
		"50001",
		"mongodb://mongo:27017")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*150)
	defer cancel()

	defer func() {
		if cnfg.MongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
