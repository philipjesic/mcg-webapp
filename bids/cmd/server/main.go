package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/bids/internal/api/routes"
	"github.com/philipjesic/mcg-webapp/bids/internal/config"
	"github.com/philipjesic/mcg-webapp/bids/internal/messaging"
	"github.com/philipjesic/mcg-webapp/bids/internal/outbox"
	"github.com/philipjesic/mcg-webapp/bids/internal/service"
	"github.com/philipjesic/mcg-webapp/bids/internal/storage"
)

func main() {
	config.LoadEnv()

	r := gin.Default()

	// TODO: Fix context objects

	// start DB
	db := storage.InitMongoClient(context.Background())

	auctionService := service.NewAuctionServiceImpl()

	// start rabbit
	msg, err := messaging.NewRabbitMQ(config.GetEnv("RABBITMQ_URI", ""), auctionService)
	if err != nil {
		panic("could not start up bids service... \nmessaging service error: " + err.Error())
	}

	// Start Outbox publisher loop
	outboxService := outbox.New(db, msg)
	outboxService.Start(context.Background(), 5*time.Second) // every 5 seconds

	routes.RegisterAPI(r, db, auctionService)

	port := config.GetEnv("PORT", "3000")

	log.Println("Listening for created bids...")
	msg.ListenForCreatedBids()

	r.Run(":" + port)

}
