# Instructions for Running the Chat Application

This application consists of two main components:
1. The web server that handles HTTP requests and WebSocket connections
2. The stock bot that processes stock quote requests

## Running with Docker Compose (Recommended)

If you have Docker and Docker Compose installed, you can run the entire application stack with a single command:

```bash
docker-compose up
```

Or using the Makefile:

```bash
make docker
```

This will start:
- PostgreSQL database
- RabbitMQ message broker
- Chat web server
- Stock bot

Access the application in your browser at: http://localhost:8080

To stop all containers:

```bash
make docker-down
```

## Manual Setup

### Prerequisites

Make sure you have the following installed:
- Go 1.18+
- PostgreSQL
- RabbitMQ

### Setup

1. Clone the repository:
```bash
git clone https://github.com/dbvitor/chat-go.git
cd chat-go
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the environment variables:
   The default configuration is in the `.env` file. Make any necessary changes to match your local environment.

4. Setup PostgreSQL:
   - Create a new database named `chatapp`
   - The application will automatically create the required tables on startup

5. Start RabbitMQ:
   Make sure RabbitMQ is running on your system.

### Running the Application

You can use the provided Makefile to run the application:

```bash
# Run both server and bot
make run

# Run only the server
make run-server

# Run only the bot
make run-bot
```

Alternatively, you can run the components manually:

1. Run the server:
```bash
go run cmd/server/main.go
```

2. Run the stock bot in a separate terminal:
```bash
go run cmd/bot/main.go
```

3. Access the application in your browser at: http://localhost:8080

## Using the Application

1. Register a new account or login if you already have one
2. Choose a chatroom to enter
3. Send messages to other users
4. Use the `/stock=<code>` command to get stock quotes (e.g., `/stock=aapl.us`)

## Building the Application

To build the application binaries:

```bash
make build
```

This will create binaries in the `build` directory.

## Testing

Run the tests with:

```bash
make test
```

Or manually:

```bash
go test ./... -v
```

## Available Make Commands

For a full list of available commands:

```bash
make help
``` 