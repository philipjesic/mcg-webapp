package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/listings/config"
	"github.com/philipjesic/mcg-webapp/listings/routes"
	storage "github.com/philipjesic/mcg-webapp/listings/storage/database"
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
