package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/handlers"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
)

func RegisterBidsHandlers(r *gin.RouterGroup, db storage.DataStore) *gin.RouterGroup {
	api := r.Group("/bids")

	bidHandler := handlers.CreateBidsHandler(db)

	api.GET("/", bidHandler.GetAll)
	api.GET("/:id", bidHandler.GetByID)
	api.POST("", bidHandler.Create)

	return api
}
