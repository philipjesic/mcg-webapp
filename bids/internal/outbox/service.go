package outbox

import (
	"context"
	"log"
	"time"

	"github.com/philipjesic/mcg-webapp/bids/internal/messaging"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OutboxService struct {
	OutboxCollection *mongo.Collection
	BidPublisher     messaging.BidPublisher
}

func New(mongoClient *storage.MongoClient, publisher messaging.BidPublisher) *OutboxService {
	collection := mongoClient.Mongo.Database("bids").Collection("outbox")
	return &OutboxService{
		OutboxCollection: collection,
		BidPublisher:     publisher,
	}
}

func (s *OutboxService) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.publishPendingBids(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *OutboxService) publishPendingBids(ctx context.Context) {
	// TODO: figure out dead head letter queue
	cursor, err := s.OutboxCollection.Find(ctx, bson.M{
		"status": bson.M{"$in": bson.A{storage.BID_OUTBOX_STATUS_PENDING, storage.BID_OUTBOX_STATUS_FAILED}},
	})
	if err != nil {
		log.Println("Outbox query failed:", err)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var outboxmsg storage.BidOutboxMessage
		if err := cursor.Decode(&outboxmsg); err != nil {
			continue
		}

		bidMsg := createBidMessage(outboxmsg)

		err = s.BidPublisher.Publish(messaging.BID_TOPIC, messaging.CREATE_BID, bidMsg)

		status := storage.BID_OUTBOX_STATUS_SENT
		if err != nil {
			log.Println("Failed to publish:", err)
			status = storage.BID_OUTBOX_STATUS_FAILED
		}

		filter := bson.M{"id": outboxmsg.ID}
		update := bson.M{
			"$set": bson.M{
				"status": status,
			},
		}
		_, _ = s.OutboxCollection.UpdateOne(ctx, filter, update)
	}
}

func createBidMessage(bid storage.BidOutboxMessage) messaging.BidMessage {
	return messaging.BidMessage{
		Event: messaging.CREATE_BID,
		Data: messaging.Bid{
			AuctionID: bid.AuctionID,
			ID:        bid.ID,
			UserID:    bid.UserID,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		},
	}
}
