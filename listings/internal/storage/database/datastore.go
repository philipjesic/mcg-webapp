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
	ID          string `bson:"id"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
}
