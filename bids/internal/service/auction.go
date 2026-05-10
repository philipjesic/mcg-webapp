package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BidEvent struct {
	Data BidEventData`json:"data"`
	Event string `json:"event"`
}

type BidEventData struct {
	AuctionID string    `json:"auction_id"`
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}


type subscription struct {
	auctionId string
	ch chan BidEvent
}

type publication struct {
	auctionId string
	bid BidEvent
}

type AuctionService interface {
	Subscribe(auctionId string) chan BidEvent
	Unsubscribe(auctionId string, bidChannel chan BidEvent)
	BroadcastBidEvent(ctx context.Context, msg amqp.Delivery) error
}

type AuctionServiceImpl struct {
	register chan subscription
	unregister chan subscription
	broadcast chan publication

	connections map[string]map[chan BidEvent]struct{}
}

func NewAuctionServiceImpl() *AuctionServiceImpl {
	a := &AuctionServiceImpl{
		register:     make(chan subscription),
		unregister:   make(chan subscription),
		broadcast:    make(chan publication),
		connections: make(map[string]map[chan BidEvent]struct{}),
	}

	go a.run()

	return a
}

func (h *AuctionServiceImpl) run() {
	for {
		select {
		// new subscriber
		case sub := <- h.register:
			if _, ok := h.connections[sub.auctionId]; !ok {
				h.connections[sub.auctionId] = make(map[chan BidEvent]struct{})
			}
			h.connections[sub.auctionId][sub.ch] = struct{}{}
		// subsriber is unregistering	
		case sub := <- h.unregister:
			subscribers, ok := h.connections[sub.auctionId]
			if !ok {
				continue
			}
			if _, exists := subscribers[sub.ch]; exists {
				delete(subscribers, sub.ch)
				close(sub.ch)
			}
			if len(subscribers) == 0 {
				delete(h.connections, sub.auctionId)
			}
		// broadcast bid to subscribers	
		case pub := <- h.broadcast:
			subscribers, ok := h.connections[pub.auctionId]
			if !ok {
				continue
			}
			for ch := range subscribers {
				select {
				case ch <- pub.bid:
				default: // subscriber was too slow to read. skip to let others read. 		
				}
			} 		 
		}
	}
}

func (h *AuctionServiceImpl) Subscribe(auctionId string) chan BidEvent {
	ch := make(chan BidEvent, 1)
	h.register <- subscription{
		auctionId: auctionId,
		ch: ch,
	}
	return ch
}

func (h *AuctionServiceImpl) Unsubscribe(auctionId string, bidChannel chan BidEvent) {
	h.unregister <- subscription{
		auctionId: auctionId,
		ch: bidChannel,
	}
}

func (h *AuctionServiceImpl) BroadcastBidEvent(ctx context.Context, msg amqp.Delivery) error {
	var bid BidEvent
	if err := json.Unmarshal(msg.Body, &bid); err != nil {
		log.Printf("invalid bid payload: %v", err)
		_ = msg.Nack(false, false) // reject and don't requeue
		return err
	}

	h.broadcast <- publication{
		auctionId: bid.Data.AuctionID,
		bid: bid,
	}
	return nil
}
