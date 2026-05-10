package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/handlers"
	"github.com/philipjesic/mcg-webapp/bids/internal/service"
)

func RegisterAuctionHandlers(r *gin.RouterGroup, auctionService service.AuctionService) *gin.RouterGroup {
	api := r.Group("/auctions")

	auctionHandler := handlers.CreateAuctionHandler(auctionService)

	api.GET("/:id", handlers.StreamHeadersMiddleware(), auctionHandler.Stream)

	return api
}
