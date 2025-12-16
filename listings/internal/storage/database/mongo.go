package database

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

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

/*
TODO: Make this function robust and performant as it will be used for the landing page.
*/
func (m *MongoClient) GetListings(ctx context.Context, ids []string) ([]Listing, error) {
	collection := m.Mongo.Database("listings").Collection("listings")
	listings := []Listing{}
	sort := bson.D{{"date_ordered", 1}}
	opts := options.Find().SetSort(sort)
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return listings, err
	}
	if err = cursor.All(ctx, &listings); err != nil {
		return listings, err
	}
	return listings, err
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

func (m *MongoClient) UpdateListingBid(ctx context.Context, auctionID, userID string, amount int, timestamp time.Time) error {
	collection := m.Mongo.Database("listings").Collection("listings")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": auctionID},
		bson.M{
			"$max": bson.M{"highestBid": amount},
			"$inc": bson.M{"bidCount": 1},
			"$set": bson.M{
				"lastBidUserId": userID,
				"lastBidAt":     timestamp,
			},
		},
	)
	return err
}
