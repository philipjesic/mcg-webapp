package storage

import (
	"context"
	"errors"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Mongo *mongo.Client
}

func InitMongoClient(ctx context.Context) *MongoClient {
	clientOptions := options.Client().
		ApplyURI(os.Getenv("MONGO_URI")).
		SetMaxPoolSize(20) // TODO: Might have to increase connection pool size later

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Could not connect to MongoDB", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ping call failed to MongoDB", err)
	}

	// TODO: fix up logging. Maybe create logging server?
	log.Println("Connected to MongoDB")

	return &MongoClient{
		Mongo: client,
	}
}

func (m *MongoClient) GetListingByID(ctx context.Context, id string) (Listing, error) {
	collection := m.Mongo.Database("listings").Collection("listings")
	listing := Listing{}
	err := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&listing)
	return listing, err
}

func (m *MongoClient) GetListings(ctx context.Context, ids []string) ([]Listing, error) {
	return []Listing{}, nil
}

func (m *MongoClient) InsertListing(ctx context.Context, l Listing) error {
	collection := m.Mongo.Database("listings").Collection("listings")
	_, err := collection.InsertOne(ctx, l)
	if err != nil {
		// TODO: add logging
		return errors.New("database error: " + err.Error())
	}
	return nil
}
