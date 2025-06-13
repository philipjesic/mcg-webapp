package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/routes"
	"github.com/philipjesic/mcg-webapp/bids/internal/config"
	"github.com/philipjesic/mcg-webapp/bids/internal/messaging"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
)

func main() {
	config.LoadEnv()

	r := gin.Default()

	// TODO: Fix context objects
	db := storage.InitMongoClient(context.Background())
	msg, err := messaging.NewRabbitMQ(config.GetEnv("RABBITMQ_URI", ""))
	if err != nil {
		panic("could not start up bids service... \nmessaging service error: " + err.Error())
	}
	routes.RegisterAPI(r, db, msg)

	port := config.GetEnv("PORT", "3000")
	r.Run(":" + port)

}
