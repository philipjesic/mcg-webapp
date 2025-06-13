package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/messaging"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
)

func RegisterAPI(r *gin.Engine, db storage.DataStore, msg messaging.BidPublisher) {
	routerGroup := r.Group("/api")

	RegisterBidsHandlers(routerGroup, db, msg)
}



