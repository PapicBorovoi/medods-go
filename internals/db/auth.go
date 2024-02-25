package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

type Auth struct {
	ID string
	RefreshToken string
}

func Create (auth *Auth) error {
	collection := DbClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("auth")

	_, err := collection.InsertOne(context.TODO(), auth)

	if err != nil {
		return err
	}

	return nil
}

func Read (id string) (*Auth, error) {
	var result *Auth
	var collection = DbClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("auth")

	err := collection.FindOne(context.Background(), bson.D{{Key: "id", Value: id}}).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func Delete (id string) error {
	var collection = DbClient.Database(os.Getenv("MONGODB_DATABASE")).Collection("auth")

	_, err := collection.DeleteOne(context.Background(), bson.D{{Key: "id", Value: id}})

	if err != nil {
		return err
	}

	return nil
}
