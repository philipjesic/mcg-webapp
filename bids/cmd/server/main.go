package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/routes"
	"github.com/philipjesic/mcg-webapp/bids/internal/config"
	"github.com/philipjesic/mcg-webapp/bids/internal/messaging"
	"github.com/philipjesic/mcg-webapp/bids/internal/outbox"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
)

func main() {
	config.LoadEnv()

	r := gin.Default()

	// TODO: Fix context objects

	// start DB
	db := storage.InitMongoClient(context.Background())

	// start rabbit
	msg, err := messaging.NewRabbitMQ(config.GetEnv("RABBITMQ_URI", ""))
	if err != nil {
		panic("could not start up bids service... \nmessaging service error: " + err.Error())
	}

	// Start Outbox publisher loop
	outboxService := outbox.New(db, msg)
	outboxService.Start(context.Background(), 5*time.Second) // every 5 seconds

	routes.RegisterAPI(r, db)

	port := config.GetEnv("PORT", "3000")
	r.Run(":" + port)

}
