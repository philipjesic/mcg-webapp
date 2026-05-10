package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/service"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
)

func RegisterAPI(r *gin.Engine, db storage.DataStore, auctionService service.AuctionService) {
	apiRouterGroup := r.Group("/api")
	serverSideEventGroup := r.Group("/api/streams")

	RegisterBidsHandlers(apiRouterGroup, db)
	RegisterAuctionHandlers(serverSideEventGroup, auctionService)
}
