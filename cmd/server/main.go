package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dbvitor/chat-go/internal/database"
	"github.com/dbvitor/chat-go/internal/handlers"
	"github.com/dbvitor/chat-go/pkg/auth"
	"github.com/dbvitor/chat-go/pkg/broker"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	err = database.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize RabbitMQ connection
	rabbitMQ, err := broker.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Consume stock results
	stockResults, err := rabbitMQ.ConsumeStockResults()
	if err != nil {
		log.Fatalf("Failed to consume stock results: %v", err)
	}

	// Initialize authentication
	auth.Initialize()

	// Create and start HTTP server
	server := handlers.NewServer(database.DB, rabbitMQ, stockResults)

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		rabbitMQ.Close()
		database.Close()
		os.Exit(0)
	}()

	// Start server
	log.Fatal(server.Start())
}
