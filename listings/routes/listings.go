package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/listings/handlers"
	storage "github.com/philipjesic/listings/storage/database"
)

func RegisterListings(r *gin.Engine, db storage.DataStore) *gin.RouterGroup {
	handler := handlers.CreateListingsHandler(db)

	api := r.Group("/listings")

	api.GET("/:id", handler.Get)
	api.POST("/", handler.Create)

	return api
}
