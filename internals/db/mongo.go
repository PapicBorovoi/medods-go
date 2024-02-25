package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DbClient *mongo.Client = nil

func Connect() error {
	fmt.Println("Database is connecting...")
	var uri = "mongodb://" + os.Getenv("MONGODB_USER") + ":" + os.Getenv("MONGODB_PASSWORD") + 
		"@localhost:27017"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	var result bson.M

	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result)

	if err != nil {
		return err
	}

	fmt.Println("Database is connected")

	DbClient = client

	return nil
}

func Close() error {
	err := DbClient.Disconnect(context.Background())

	if err != nil {
		return err
	}

	DbClient = nil

	return nil
}