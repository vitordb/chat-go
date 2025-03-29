package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dbvitor/chat-go/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents the RabbitMQ connection
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ connection
func NewRabbitMQ() (*RabbitMQ, error) {
	// Get connection details from environment variables
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")

	// Create connection string
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}

	// Create channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Create queues if they don't exist
	stockQueue := os.Getenv("RABBITMQ_STOCK_QUEUE")
	resultQueue := os.Getenv("RABBITMQ_RESULT_QUEUE")

	// Declare stock queue
	_, err = channel.QueueDeclare(
		stockQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	// Declare result queue
	_, err = channel.QueueDeclare(
		resultQueue, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// PublishStockRequest publishes a stock request to the stock queue
func (r *RabbitMQ) PublishStockRequest(chatroomID, stockCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create request message
	request := map[string]string{
		"chatroom_id": chatroomID,
		"stock_code":  stockCode,
	}

	// Marshal request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Publish message to stock queue
	return r.channel.PublishWithContext(
		ctx,
		"",                                // exchange
		os.Getenv("RABBITMQ_STOCK_QUEUE"), // routing key
		false,                             // mandatory
		false,                             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// PublishStockResult publishes a stock result to the result queue
func (r *RabbitMQ) PublishStockResult(response *models.StockResponse, chatroomID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create result message
	result := map[string]interface{}{
		"chatroom_id": chatroomID,
		"symbol":      response.Symbol,
		"price":       response.Price,
		"error":       response.Error,
	}

	// Marshal result to JSON
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Publish message to result queue
	return r.channel.PublishWithContext(
		ctx,
		"",                                 // exchange
		os.Getenv("RABBITMQ_RESULT_QUEUE"), // routing key
		false,                              // mandatory
		false,                              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// ConsumeStockRequests consumes stock requests from the stock queue
func (r *RabbitMQ) ConsumeStockRequests() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		os.Getenv("RABBITMQ_STOCK_QUEUE"), // queue
		"",                                // consumer
		true,                              // auto-ack
		false,                             // exclusive
		false,                             // no-local
		false,                             // no-wait
		nil,                               // args
	)
}

// ConsumeStockResults consumes stock results from the result queue
func (r *RabbitMQ) ConsumeStockResults() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		os.Getenv("RABBITMQ_RESULT_QUEUE"), // queue
		"",                                 // consumer
		true,                               // auto-ack
		false,                              // exclusive
		false,                              // no-local
		false,                              // no-wait
		nil,                                // args
	)
}
