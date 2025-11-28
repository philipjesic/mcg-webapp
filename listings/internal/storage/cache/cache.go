package cache

import ("context")

type Cache interface {
	UpdateBid(ctx context.Context, auctionID, userID string, amount int) error
}

