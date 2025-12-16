package main

import (
	"context"
	"strconv"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/philipjesic/mcg-webapp/listings/internal/api/routes"
	"github.com/philipjesic/mcg-webapp/listings/internal/config"
	"github.com/philipjesic/mcg-webapp/listings/internal/messaging"
	"github.com/philipjesic/mcg-webapp/listings/internal/storage/bids"
	"github.com/philipjesic/mcg-webapp/listings/internal/storage/buffer"
	"github.com/philipjesic/mcg-webapp/listings/internal/storage/cache"
	storage "github.com/philipjesic/mcg-webapp/listings/internal/storage/database"
)

func main() {
	config.LoadEnv()

	r := gin.Default()

	// TODO: Fix context objects
	db := storage.InitMongoClient(context.Background())
	log.Println("Connected to MongoDB...")

	redis, cacheErr := cache.NewRedisCache(config.GetEnv("REDIS_URI", ""))
	if cacheErr != nil {
		panic(cacheErr)
	}
	log.Println("Connected to Redis...")

	flushInterval, _ := strconv.Atoi(config.GetEnv("BUFFER_FLUSH_INTERVAL", "5"))
	bidBuffer := buffer.NewBuffer(db, time.Duration(flushInterval)*time.Second)

	bidHandler := bids.NewHandler(bidBuffer, redis)

	rabbbitMQ, rabbitErr := messaging.NewRabbitMQ(config.GetEnv("RABBITMQ_URI", ""), bidHandler)
	if rabbitErr != nil {
		panic(rabbitErr)
	}
	log.Println("Connected to RabbitMQ...")

	rabbbitMQ.ListenForCreatedBids()
	
	routes.RegisterAPI(r, db)

	port := config.GetEnv("PORT", "3000")
	r.Run(":" + port)
}
