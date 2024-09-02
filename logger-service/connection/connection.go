package connection

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongo(url string) (*mongo.Client, error) {
	clientOption := options.Client().ApplyURI(url)
	clientOption.SetAuth(options.Credential{
		Username: "root",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	return c, nil
}
