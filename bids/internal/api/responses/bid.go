package responses

import "time"

type BidsResponse struct {
	Data []BidResponseBody `json:"data"`
}

type BidResponseBody struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	AuctionID string    `json:"auction_id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount" `
	Timestamp time.Time `json:"timestamp"`
}
