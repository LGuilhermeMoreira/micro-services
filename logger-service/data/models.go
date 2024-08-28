package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

type Model struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mc *mongo.Client) Model {
	mongoClient = mc
	return Model{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := mongoClient.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), entry)

	if err != nil {
		log.Println("Error inserting into logs", err)
	}

	return err
}

func (l *LogEntry) GetAll() ([]LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	collection := mongoClient.Database("logs").Collection("logs")
	options := options.Find()
	options.SetSort(bson.D{{"create_at", -1}})
	response, err := collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		log.Println("Error finding all logs", err)
		return nil, err
	}
	defer response.Close(ctx)

	var itens []LogEntry
	for response.Next(ctx) {
		var item LogEntry

		err := response.Decode(&item)

		if err != nil {
			log.Println("Error decoding log into slice", err)
			return nil, err
		}

		itens = append(itens, item)
	}

	return itens, nil
}
