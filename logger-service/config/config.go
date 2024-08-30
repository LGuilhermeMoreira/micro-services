package config

import (
	"logger/connection"
	"logger/data"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	WebPort     string
	RpcPort     string
	MongoClient *mongo.Client
	GrpcPort    string
	Models      data.Model
}

func NewConfig(webport, rpcport, grpcport, mongourl string) (*Config, error) {
	client, err := connection.ConnectToMongo(mongourl)
	if err != nil {
		return nil, err
	}
	return &Config{
		WebPort:     webport,
		RpcPort:     rpcport,
		GrpcPort:    grpcport,
		MongoClient: client,
		Models:      data.New(client),
	}, nil
}
