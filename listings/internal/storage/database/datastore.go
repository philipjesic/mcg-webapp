package storage

import "context"

type DataStore interface {
	GetListingByID(ctx context.Context, id string) (Listing, error)
	GetListings(ctx context.Context, ids []string) ([]Listing, error)
	InsertListing(ctx context.Context, l Listing) error
}

type Listing struct {
	ID          string `bson:"id"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
}
