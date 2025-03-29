# Chat Application with Stock Bot

A simple browser-based chat application using Go, with stock quote functionality.

## Features

- User registration and login
- Real-time chat in chatrooms
- Stock quote command `/stock=stock_code`
- Message broker integration with RabbitMQ
- Last 50 messages displayed, ordered by timestamp

## Architecture

The application is built using a clean architecture with the following components:

- Web server for handling HTTP requests and WebSocket connections
- Stock bot for fetching stock quotes from external API
- RabbitMQ for message broker communication
- Database for storing users, messages, and chatrooms

## Setup

### Prerequisites

- Go 1.18+
- RabbitMQ
- PostgreSQL

### Running the Application

1. Clone the repository
2. Configure the environment variables in config file
3. Start the RabbitMQ server
4. Start the PostgreSQL server
5. Run the server: `go run cmd/server/main.go`
6. Run the bot: `go run cmd/bot/main.go`
7. Access the application at http://localhost:8080

## Usage

- Register and login to access chat
- Enter messages in the chatroom
- Use `/stock=stock_code` command to get stock quotes (e.g., `/stock=aapl.us`) 