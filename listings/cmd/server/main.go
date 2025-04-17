package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/listings/config"
	"github.com/philipjesic/listings/routes"
	storage "github.com/philipjesic/listings/storage/database"
)

func main() {
	config.LoadEnv()

	r := gin.Default()

	// TODO: Fix context objects
	db := storage.InitMongoClient(context.TODO())
	routes.RegisterRoutes(r, db)

	port := config.GetEnv("PORT", "3000")
	r.Run(":" + port)
}
