package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/listings/internal/config"
	"github.com/philipjesic/mcg-webapp/listings/internal/api/routes"
	storage "github.com/philipjesic/mcg-webapp/listings/internal/storage/database"
)

func main() {
	config.LoadEnv()

	r := gin.Default()

	// TODO: Fix context objects
	db := storage.InitMongoClient(context.Background())
	routes.RegisterAPI(r, db)

	port := config.GetEnv("PORT", "3000")
	r.Run(":" + port)
}
