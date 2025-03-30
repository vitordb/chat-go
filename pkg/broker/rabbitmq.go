package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dbvitor/chat-go/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Wrapper for RabbitMQ connection
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connects to RabbitMQ and configures the necessary queues
func NewRabbitMQ() (*RabbitMQ, error) {
	// Gets configurations from .env or uses default values
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672"
	}

	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}

	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		password = "guest"
	}

	// Log to help with debugging
	log.Printf("Connecting to RabbitMQ at amqp://%s:***@%s:%s/", user, host, port)

	// Builds connection URL
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	// Connects
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Creates channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Queue names (from .env or uses default)
	stockQueue := os.Getenv("RABBITMQ_STOCK_QUEUE")
	if stockQueue == "" {
		stockQueue = "stock_requests"
	}

	resultQueue := os.Getenv("RABBITMQ_RESULT_QUEUE")
	if resultQueue == "" {
		resultQueue = "stock_results"
	}

	// Creates queue for stock requests
	_, err = channel.QueueDeclare(
		stockQueue, // name
		true,       // durable (survives restart)
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

	// Creates queue for results
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

// Closes connection and releases resources
func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// Sends stock request to the bot for processing
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
		log.Printf("ERROR marshaling stock request: %v", err)
		return err
	}

	// Get queue name from environment or use default
	queueName := os.Getenv("RABBITMQ_STOCK_QUEUE")
	if queueName == "" {
		queueName = "stock_requests" // Default queue name
	}

	log.Printf("Publishing stock request for %s to queue %s", stockCode, queueName)

	// Publish message to stock queue
	err = r.channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)

	if err != nil {
		log.Printf("ERROR publishing stock request: %v", err)
		return err
	}

	log.Printf("Successfully published stock request for %s to queue %s", stockCode, queueName)
	return nil
}

// Sends stock quote response back to the server
func (r *RabbitMQ) PublishStockResult(response *models.StockResponse, chatroomID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Data to send
	result := map[string]interface{}{
		"chatroom_id": chatroomID,
		"symbol":      response.Symbol,
		"price":       response.Price,
		"error":       response.Error,
	}

	// Converts to JSON
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// Publishes to the results queue
	return r.channel.PublishWithContext(
		ctx,
		"",                                 // exchange
		os.Getenv("RABBITMQ_RESULT_QUEUE"), // destination queue
		false,                              // mandatory
		false,                              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// Receives stock requests (used by the bot)
func (r *RabbitMQ) ConsumeStockRequests() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		os.Getenv("RABBITMQ_STOCK_QUEUE"), // queue
		"",                                // consumer (empty generates unique ID)
		true,                              // auto-ack
		false,                             // exclusive
		false,                             // no-local
		false,                             // no-wait
		nil,                               // args
	)
}

// Receives stock responses (used by the server)
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
