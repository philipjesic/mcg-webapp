package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/philipjesic/mcg-webapp/bids/internal/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	connection       *amqp.Connection
	channel          *amqp.Channel
	BidCreationQueue string
	AuctionService   service.AuctionService
}

func NewRabbitMQ(url string, auctionService service.AuctionService) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Panicf("%s: %s", "Failed to create a channel to RabbitMQ", err)
		return nil, err
	}

	// Initialize the exchange for bids
	err = ch.ExchangeDeclare(
		BID_TOPIC, // name
		"topic",   // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		failOnError(err, "Failed to declare BID exchange in RabbitMQ", ch, conn)
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		failOnError(err, "Failed to declare QUEUE in RabbitMQ", ch, conn)
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		CREATE_BID,
		BID_TOPIC,
		false,
		nil,
	)

	if err != nil {
		failOnError(err, "Failed to bind Queue to exchange in RabbitMQ", ch, conn)
		return nil, err
	}

	return &RabbitMQ{
		connection:       conn,
		channel:          ch,
		BidCreationQueue: q.Name,
		AuctionService:   auctionService,
	}, nil
}

func (r *RabbitMQ) Publish(topic, key string, bidMsg BidMessage) error {
	body, err := json.Marshal(bidMsg)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		topic, // exchange
		key,   // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) ListenForCreatedBids() {
	msgs, err := r.channel.Consume(
		r.BidCreationQueue, // queue
		"",                 // consumer
		true,               // auto ack
		false,              // exclusive
		false,              // no local
		false,              // no wait
		nil,                // args
	)

	if err != nil {
		failOnError(err, "Failed to bind Queue to exchange in RabbitMQ", r.channel, r.connection)
	}

	go func() {
		for msg := range msgs {
			msgCopy := msg
			err := r.AuctionService.BroadcastBidEvent(context.Background(), msgCopy)
			if err != nil {
				log.Printf("error handling created bid event: %v", err)
			}
		}
	}()
}

func failOnError(err error, msg string, ch *amqp.Channel, conn *amqp.Connection) {
	ch.Close()
	conn.Close()
	log.Panicf("%s: %s", msg, err)
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.connection != nil {
		r.connection.Close()
	}
}
