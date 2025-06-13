package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
	"github.com/philipjesic/mcg-webapp/bids/internal/messaging"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/handlers"
)

func RegisterBidsHandlers(r *gin.RouterGroup, db storage.DataStore, msg messaging.BidPublisher) *gin.RouterGroup {
	api := r.Group("/bids")
	
	bidHandler := handlers.CreateBidsHandler(db, msg)

	api.GET("/", bidHandler.GetAll)
	api.GET("/:id", bidHandler.GetByID)
	api.POST("", bidHandler.Create)

	return api
}