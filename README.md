# Chat Application with Stock Bot

A browser-based chat application with stock quote functionality.

## Features

- User registration and login
- Real-time chat
- Stock quote command `/stock=stock_code` (e.g., `/stock=aapl.us`)
- Message broker integration with RabbitMQ
- Last 50 messages displayed, ordered by timestamp

## Running the Application

### Prerequisites

- Go 1.18+
- Docker
- Docker Compose

### Quick Start (Recommended)

To start all services (PostgreSQL, RabbitMQ, web server, and stock bot) with a single command:

```bash
make run-all
```

This command will:
1. Stop any previous Go processes
2. Start the web server and bot in background
3. Show instructions to access the application

Access the application at: http://localhost:8080

### First-Time Setup or Reset

If you need to start with a fresh database or are having issues, use:

```bash
make reset-all
```

This command will:
1. Stop any previous instances
2. Restart Docker containers (PostgreSQL and RabbitMQ)
3. Reset the database
4. Start the web server and bot in background
5. Show instructions to access the application

### Viewing Logs

To see the logs:
```bash
# Server logs
tail -f server.log

# Bot logs
tail -f bot.log
```

### Stopping All Services

```bash
make stop-all
```

## Testing the Stock Quote Functionality

1. Access http://localhost:8080
2. Register a new user
3. Enter the chat room
4. Type the command `/stock=aapl.us` to get Apple's stock quote

## Architecture

The application is built using a clean architecture with the following components:

- Web server for handling HTTP requests and WebSocket connections
- Stock bot for fetching stock quotes from external API
- RabbitMQ for service communication
- PostgreSQL database for storing users, messages, and chat rooms

## Usage

- Register and login to access chat
- Enter messages in the chatroom
- Use `/stock=stock_code` command to get stock quotes (e.g., `/stock=aapl.us`) 