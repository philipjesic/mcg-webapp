package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/listings/handlers"
	storage "github.com/philipjesic/mcg-webapp/listings/storage/database"
)

func RegisterListingHandlers(r *gin.RouterGroup, db storage.DataStore) *gin.RouterGroup {
	handler := handlers.CreateListingsHandler(db)

	api := r.Group("/listings")

	api.GET("/:id", handler.Get)
	api.POST("/", handler.Create)

	return api
}
