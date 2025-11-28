package bids

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/philipjesic/mcg-webapp/listings/internal/storage/buffer"
	"github.com/philipjesic/mcg-webapp/listings/internal/storage/cache"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	buffer *buffer.Buffer
	cache  cache.Cache
}

func NewHandler(buffer *buffer.Buffer, cache cache.Cache) *Handler {
	return &Handler{
		buffer: buffer,
		cache: cache,
	}
}

type BidEvent struct {
	Data struct {
		AuctionID string    `json:"auction_id"`
		ID        string    `json:"id" `
		UserID    string    `json:"user_id"`
		Amount    int       `json:"amount" `
		Timestamp time.Time `json:"timestamp"`
	} `json:"data"`
	Event string `json:"event"`
}

func (h *Handler) HandleCreateBidEvent(ctx context.Context, msg *amqp.Delivery) error {
	// 1. Parse the incoming message
	var bid BidEvent
	if err := json.Unmarshal(msg.Body, &bid); err != nil {
		log.Printf("invalid bid payload: %v", err)
		_ = msg.Nack(false, false) // reject and don't requeue
		return err
	}

	// 2. Write to Redis (fast path)
	if err := h.cache.UpdateBid(ctx, bid.Data.AuctionID, bid.Data.UserID, bid.Data.Amount); err != nil {
		log.Printf("failed to update Redis: %v", err)
		_ = msg.Nack(false, true) // requeue to retry later
		return err
	}

	// 3. Add to buffer (Mongo will be flushed async)
	h.buffer.Add(createBufferEntry(bid, msg))
	return nil
}

func createBufferEntry(bidEvent BidEvent, msg *amqp.Delivery) buffer.Entry {
	return buffer.Entry{
		AuctionID: bidEvent.Data.AuctionID,
		ID: bidEvent.Data.ID,
		UserID: bidEvent.Data.UserID,
		Amount: bidEvent.Data.Amount,
		Timestamp: bidEvent.Data.Timestamp,
		Msg: msg,
	}
}
