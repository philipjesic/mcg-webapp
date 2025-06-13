package storage

import (
	"context"
	"time"
)

type DataStore interface {
	GetBidByID(ctx context.Context, id string) (Bid, error)
	GetBids(ctx context.Context, ids []string) ([]Bid, error)
	InsertBid(ctx context.Context, bid Bid) error
}

type Bid struct {
	AuctionID string    `bson:"auction_id"`
	ID        string    `bson:"id" `
	UserID    string    `bson:"user_id"`
	Amount    int       `bson:"amount" `
	Timestamp time.Time `bson:"timestamp"`
}
