package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/philipjesic/listings/storage/database"
)

func RegisterRoutes(r *gin.Engine, db storage.DataStore) {
	RegisterListings(r, db)
}

