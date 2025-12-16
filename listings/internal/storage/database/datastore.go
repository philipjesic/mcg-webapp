package database

import (
	"context"
	"time"
)

type DataStore interface {
	GetListingByID(ctx context.Context, id string) (Listing, error)
	GetListings(ctx context.Context, ids []string) ([]Listing, error)
	InsertListing(ctx context.Context, l Listing) error
	UpdateListingBid(ctx context.Context, auctionID, userID string, amount int, timestamp time.Time) error
}

type Listing struct {
	ID             string         `bson:"id"`
	Title          string         `bson:"title"`
	Description    string         `bson:"description"`
	Location       string         `bson: "location"`
	Endtime        string         `bson: "endtime"`
	Seller         Seller         `bson: "seller`
	Specifications Specifications `bson: "specifications"`
}

type Seller struct {
	Name     string `bson:"name"`
	Location string `bson:"location"`
	ID       string `bson:"id"`
}

type Specifications struct {
	Make         string `bson:"make"`
	Model        string `bson:"model"`
	Engine       string `bson:"engine"`
	Colour       string `bson:"colour"`
	Transmission string `bson:"transmission"`
}
