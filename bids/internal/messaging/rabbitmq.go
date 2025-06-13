package messaging

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
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
		ch.Close()
		conn.Close()
		log.Panicf("%s: %s", "Failed to declare BID exchange in RabbitMQ", err)
		return nil, err
	}

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
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

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.connection != nil {
		r.connection.Close()
	}
}
