package handlers

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/service"
)

type Auction struct {
	auctionService service.AuctionService
}

func CreateAuctionHandler(auctionService service.AuctionService) *Auction {
	return &Auction{
		auctionService: auctionService,
	}
}

func (a *Auction) Stream(ctx *gin.Context) {
	auctionId := ctx.Param("id")
	bidChannel := a.auctionService.Subscribe(auctionId)
	defer a.auctionService.Unsubscribe(auctionId, bidChannel)

	ctx.SSEvent("Connected", "Listening for auction bids...")

	ctx.Stream(func(w io.Writer) bool {
		select {
		case bid, ok := <-bidChannel:
			if !ok {
				ctx.SSEvent("Disconnect", "The auction is now closed")
				return false
			}

			ctx.SSEvent("Bid", bid)
			return true
		case <-ctx.Request.Context().Done():
			return false
		}
	})
}

func StreamHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
