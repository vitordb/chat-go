# Instructions for Running the Chat Application

This application consists of two main components:
1. The web server that handles HTTP requests and WebSocket connections
2. The stock bot that processes stock quote requests

## Prerequisites

Make sure you have the following installed:
- Docker
- Docker Compose
- Go 1.18+ (for local execution without Docker)

## Quick Start 

The recommended way to start the application is:

```bash
make run-all
```

This command will:
- Stop any previous Go processes
- Start the web server and stock bot in background
- Show instructions to access the application

If you are starting the application for the first time or need to start with a clean database:

```bash
make reset-all
```

This command will:
- Stop all existing services
- Remove Docker containers and volumes (resetting the database)
- Restart Docker containers (PostgreSQL and RabbitMQ)
- Start the web server and stock bot in background

## Different Ways to Run the Application

### Running Everything in Background (Recommended for Development)

```bash
# Start the application with existing Docker containers
make run-all

# OR, for a clean start with fresh database
make reset-all
```

### Running with Logs in Foreground

```bash
# Start Docker containers, server and bot with logs in the terminal
make run
```

### Running Just the Docker Containers

```bash
# Start only the Docker containers (PostgreSQL and RabbitMQ)
make docker
```

### Stopping All Services

```bash
make stop-all
```

## Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/dbvitor/chat-go.git
cd chat-go
```

2. Configure the `.env` file:
```bash
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chatapp
DB_SSL_MODE=disable

# RabbitMQ Configuration
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# Server Configuration
SERVER_PORT=8080
```

## Viewing Logs

```bash
# Server logs
tail -f server.log

# Bot logs
tail -f bot.log
```

## Testing the Application

1. Access the application in your browser:
```
http://localhost:8080
```

2. Register a new user
3. Enter a chat room
4. Try the stock quote command:
```
/stock=aapl.us
```

