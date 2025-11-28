package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) (*Redis, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	return &Redis{client: client}, nil
}

func (r *Redis) UpdateBid(ctx context.Context, auctionID, userID string, amount int) error {
	key := fmt.Sprintf("auction:%s", auctionID)

	// Atomically update highestBid and increment bidCount
	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		// Set lastBidUserId and lastBidAt outside the transaction
		_, err := tx.HSet(ctx, key, map[string]interface{}{
			"lastBidUserId": userID,
			"lastBidAt":     time.Now().Format(time.RFC3339),
		}).Result()
		if err != nil {
			return err
		}

		currentBidStr, err := tx.HGet(ctx, key, "highestBid").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		currentBid, _ := strconv.Atoi(currentBidStr)

		pipe := tx.Pipeline()
		if amount > currentBid {
			pipe.HSet(ctx, key, "highestBid", amount)
		}
		pipe.HIncrBy(ctx, key, "bidCount", 1)
		_, err = pipe.Exec(ctx)
		return err
	}, key)

	return err
}

func (r *Redis) GetLiveBid(ctx context.Context, auctionID string) (map[string]string, error) {
	key := fmt.Sprintf("auction:%s", auctionID)
	return r.client.HGetAll(ctx, key).Result()
}
