package messaging

import (
	"context"
	"log"

	"github.com/philipjesic/mcg-webapp/listings/internal/storage/bids"
	amqp "github.com/rabbitmq/amqp091-go"
)

const CREATE_BID = "bid.create"

const BID_TOPIC = "bids"

type RabbitMQ struct {
	connection       *amqp.Connection
	channel          *amqp.Channel
	BidHandler       *bids.Handler
	BidCreationQueue string
}

func NewRabbitMQ(url string, bidHandler *bids.Handler) (*RabbitMQ, error) {
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
		"",    // let RabbitMQ generate a random queue name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false,
		nil,
	)

	if err != nil {
		failOnError(err, "Failed to declare QUEUE exchange in RabbitMQ", ch, conn)
		return nil, err
	}

	// Bind to BIDS topic
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
		BidHandler:       bidHandler,
	}, nil
}

func (r *RabbitMQ) ListenForCreatedBids() {
	msgs, err := r.channel.Consume(
		r.BidCreationQueue,
		"",
		false, // <-- autoAck false, we want to ack manually
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		failOnError(err, "Failed to bind Queue to exchange in RabbitMQ", r.channel, r.connection)
	}

	go func() {
		for d := range msgs {
			msgCopy := d
			log.Printf("Received a message: %s", d.Body)
			ctx := context.Background()
			err := r.BidHandler.HandleCreateBidEvent(ctx, &msgCopy)
			if err != nil {
				log.Printf("Create Bid event persistence failed: %v", err)
				// DO NOT ACK - message will be requeued or dead-lettered
				continue
			}
		}
	}()
}

func failOnError(err error, msg string, ch *amqp.Channel, conn *amqp.Connection) {
	ch.Close()
	conn.Close()
	log.Panicf("%s: %s", msg, err)
}
