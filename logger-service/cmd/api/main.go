package main

import "logger/cmd/config"

func main() {
	_, err := config.NewConfig("80",
		"5001",
		"50001",
		"mongodb://mongo:27017")
	if err != nil {
		panic(err)
	}
}
