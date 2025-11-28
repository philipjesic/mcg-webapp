package buffer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/philipjesic/mcg-webapp/listings/internal/storage/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Buffer struct {
	mu       sync.Mutex
	entries  map[string][]Entry
	db       database.DataStore
	interval time.Duration
}

type Entry struct {
	AuctionID string
	ID        string
	UserID    string
	Amount    int
	Timestamp time.Time
	Msg       *amqp.Delivery
}

func NewBuffer(db database.DataStore, flushInterval time.Duration) *Buffer {
	b := &Buffer{
		entries:  make(map[string][]Entry),
		db:       db,
		interval: flushInterval,
	}
	go b.runFlushLoop()
	return b
}

func (b *Buffer) Add(entry Entry) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[entry.AuctionID] = append(b.entries[entry.AuctionID], entry)
}

func (b *Buffer) runFlushLoop() {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	for range ticker.C {
		b.flush()
	}
}

func (b *Buffer) flush() {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, entries := range b.entries {
		for _, entry := range entries {
			err := b.db.UpdateListingBid(ctx, entry.AuctionID, entry.UserID, entry.Amount, entry.Timestamp)
			if err != nil {
				log.Println("error saving bid buffer entries to database: " + err.Error())
				continue
			} else {
				entry.Msg.Ack(true)
			}
		}
	}

	// Clear buffer after flush
	b.entries = make(map[string][]Entry)
}
