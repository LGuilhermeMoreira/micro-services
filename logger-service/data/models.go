package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	opt := options.Find()
	opt.SetSort(bson.D{{"create_at", -1}})
	response, err := collection.Find(context.TODO(), bson.D{}, opt)
	if err != nil {
		log.Println("Error finding all logs", err)
		return nil, err
	}
	defer response.Close(ctx)

	var logs []LogEntry
	for response.Next(ctx) {
		var item LogEntry

		err := response.Decode(&item)

		if err != nil {
			log.Println("Error decoding log into slice", err)
			return nil, err
		}

		logs = append(logs, item)
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := mongoClient.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println("Error finding log", err)
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.D{{"_id", docId}}).Decode(&entry)

	if err != nil {
		log.Println("Error decoding entry", err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	collection := mongoClient.Database("logs").Collection(name)
	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping log", err)
	}
	return err
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := mongoClient.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		log.Println("Error finding log", err)
		return nil, err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": docId}, bson.D{{"$set", bson.D{
		{"name", l.Name},
		{"data", l.Data},
		{"created_at", l.CreatedAt},
		{"updated_at", time.Now()},
	}}})

	if err != nil {
		log.Println("Error updating log", err)
		return nil, err
	}

	return result, nil
}
