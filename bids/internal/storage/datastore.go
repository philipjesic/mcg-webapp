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

type BidOutboxMessage struct {
	AuctionID string    `bson:"auction_id"`
	ID        string    `bson:"id" `
	UserID    string    `bson:"user_id"`
	Amount    int       `bson:"amount" `
	Timestamp time.Time `bson:"timestamp"`
	Event     string    `bson:"event_type"`
	Status    string    `bson:"status"`
}

const BID_OUTBOX_STATUS_PENDING = "pending"
const BID_OUTBOX_STATUS_SENT = "sent"
const BID_OUTBOX_STATUS_FAILED = "failed"

func createBidOutboxMessage(bid Bid, event string, status string) BidOutboxMessage {
	return BidOutboxMessage{
		Event:     event,
		AuctionID: bid.AuctionID,
		ID:        bid.ID,
		UserID:    bid.UserID,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
		Status:    status,
	}
}
