package routes

import (
	"github.com/gin-gonic/gin"
	storage "github.com/philipjesic/mcg-webapp/listings/internal/storage/database"
)

func RegisterAPI(r *gin.Engine, db storage.DataStore) {

	routerGroup := r.Group("/api")

	RegisterListingHandlers(routerGroup, db)
}
