package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dbvitor/chat-go/internal/models"
	"github.com/dbvitor/chat-go/internal/services"
	"github.com/dbvitor/chat-go/pkg/broker"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Check STOCK_API_URL variable
	stockApiUrl := os.Getenv("STOCK_API_URL")
	if stockApiUrl == "" {
		log.Println("ERROR: STOCK_API_URL environment variable not configured!")
	} else {
		log.Printf("STOCK_API_URL configured: %s", stockApiUrl)
	}

	// Initialize RabbitMQ connection
	rabbitMQ, err := broker.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Create stock service
	stockService := services.NewStockService()

	// Consume stock requests
	stockRequests, err := rabbitMQ.ConsumeStockRequests()
	if err != nil {
		log.Fatalf("Failed to consume stock requests: %v", err)
	}

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		rabbitMQ.Close()
		os.Exit(0)
	}()

	log.Println("Stock bot started. Waiting for requests...")

	// Process stock requests
	for delivery := range stockRequests {
		go func(delivery []byte) {
			// Parse request
			var request map[string]string
			if err := json.Unmarshal(delivery, &request); err != nil {
				log.Printf("Error parsing request: %v", err)
				return
			}

			// Extract data
			chatroomID := request["chatroom_id"]
			stockCode := request["stock_code"]

			log.Printf("Processing stock request: %s for chatroom %s", stockCode, chatroomID)

			// Get stock quote
			log.Printf("Getting quote for %s", stockCode)
			stockResponse, err := stockService.GetStockQuote(stockCode)
			if err != nil {
				log.Printf("Error getting stock quote: %v", err)
				stockResponse = &models.StockResponse{
					Symbol: stockCode,
					Error:  "internal server error",
				}
			}

			if stockResponse.Error != "" {
				log.Printf("Error retrieving quote: %s", stockResponse.Error)
			} else {
				log.Printf("Quote successfully retrieved: %s = $%.2f", stockResponse.Symbol, stockResponse.Price)
			}

			// Publish result
			log.Printf("Publishing result to chatroom %s", chatroomID)
			err = rabbitMQ.PublishStockResult(stockResponse, chatroomID)
			if err != nil {
				log.Printf("Error publishing result: %v", err)
			} else {
				log.Printf("Result successfully published to chatroom %s", chatroomID)
			}
		}(delivery.Body)
	}
}
