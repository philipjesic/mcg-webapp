package messaging

import "time"

const CREATE_BID = "bid.create"

const BID_TOPIC = "bids"

type BidPublisher interface {
	Publish(topic, key string, bidMsg BidMessage) error
}

type Bid struct {
	AuctionID string    `json:"auction_id"`
	ID        string    `json:"id" `
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount" `
	Timestamp time.Time `json:"timestamp"`
}

type BidMessage struct {
	Bid Bid `json:"data"`
}
