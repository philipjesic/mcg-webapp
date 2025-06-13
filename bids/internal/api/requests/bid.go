package requests

import "time"

type BidCreateRequestBody struct {
	AuctionID string    `json:"auction_id" binding:"required"`
	UserID    string    `json:"user_id" binding:"required"`
	Amount    int       `json:"amount" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

type BidCreateRequest struct {
	Data BidCreateRequestBody `json:"data" binding:"required"`
}
