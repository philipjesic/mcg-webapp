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

func (m *MongoClient) GetBidByID(ctx context.Context, id string) (Bid, error) {
	collection := m.Mongo.Database("bids").Collection("bids")
	bid := Bid{}
	err := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&bid)
	return bid, err
}

func (m *MongoClient) GetBids(ctx context.Context, ids []string) ([]Bid, error) {
	collection := m.Mongo.Database("bids").Collection("bids")
	bids := []Bid{}
	sort := bson.D{{Key: "date", Value: 1}}
	opts := options.Find().SetSort(sort)
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return bids, err
	}
	if err = cursor.All(ctx, &bids); err != nil {
		return bids, err
	}
	return bids, err
}

func (m *MongoClient) InsertBid(ctx context.Context, b Bid) error {
	session, err := m.Mongo.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		// Insert bid
		collection := m.Mongo.Database("bids").Collection("bids")
		_, err := collection.InsertOne(sc, b)
		if err != nil {
			// TODO: add logging
			return nil, errors.New("database error: " + err.Error())
		}

		// Create outbox message
		outboxMsg := createBidOutboxMessage(b, "bid.create", BID_OUTBOX_STATUS_PENDING)
		outBoxCollection := m.Mongo.Database("bids").Collection("outbox")
		_, err = outBoxCollection.InsertOne(sc, outboxMsg)
		if err != nil {
			return nil, errors.New("database error: " + err.Error())
		}
		return nil, err
	})

	return err

}
