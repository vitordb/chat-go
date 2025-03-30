# Makefile for the Chat-Go project

# Variables
GO=go
DOCKER=docker compose
SERVER_CMD=./cmd/server/main.go
BOT_CMD=./cmd/bot/main.go
BUILD_DIR=./build

# Colors for better readability
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
RESET=\033[0m

# Main commands
.PHONY: all run run-all run-server run-bot build clean test docker docker-down help reset-all

all: build

# Start the entire application with Docker Compose
docker:
	@echo "$(GREEN)Starting application with Docker Compose...$(RESET)"
	@$(DOCKER) up -d
	@echo "$(GREEN)Services started! Access http://localhost:8080$(RESET)"

# Stop Docker containers
docker-down:
	@echo "$(YELLOW)Stopping all containers...$(RESET)"
	@$(DOCKER) down
	@echo "$(GREEN)Containers stopped successfully.$(RESET)"

# Run all services in background mode (main command to start the application)
run-all:
	@echo "$(GREEN)Stopping any previous Go processes...$(RESET)"
	@pkill -f "go run" || true
	@echo "$(YELLOW)Checking if port 8080 is in use...$(RESET)"
	@lsof -i :8080 -t | xargs kill -9 2>/dev/null || true
	@echo "$(GREEN)Starting server in background...$(RESET)"
	@$(GO) run $(SERVER_CMD) > server.log 2>&1 &
	@echo "$(YELLOW)Waiting for server to initialize (3 seconds)...$(RESET)"
	@sleep 3
	@echo "$(GREEN)Starting stock bot in background...$(RESET)"
	@$(GO) run $(BOT_CMD) > bot.log 2>&1 &
	@echo "$(GREEN)All services started!$(RESET)"
	@echo "$(GREEN)Access the application at: http://localhost:8080$(RESET)"
	@echo "$(YELLOW)View server logs: tail -f server.log$(RESET)"
	@echo "$(YELLOW)View bot logs: tail -f bot.log$(RESET)"
	@echo "$(YELLOW)Stop all services: make stop-all$(RESET)"
	@echo "$(YELLOW)Checking service status...$(RESET)"
	@sleep 2
	@if lsof -i :8080 | grep LISTEN >/dev/null; then \
		echo "$(GREEN)Server running on port 8080$(RESET)"; \
	else \
		echo "$(RED)Warning: Server not running on port 8080, check server.log for errors$(RESET)"; \
	fi
	@echo "$(GREEN)Done! Try using /stock=aapl.us in the chat$(RESET)"

# Reset and restart everything, including database (use when you need a fresh start)
reset-all:
	@echo "$(YELLOW)Stopping all services...$(RESET)"
	@pkill -f "go run" || true
	@$(DOCKER) down -v
	@echo "$(GREEN)Removing logs...$(RESET)"
	@rm -f *.log
	@echo "$(GREEN)Starting Docker containers with fresh database...$(RESET)"
	@$(DOCKER) up -d
	@echo "$(YELLOW)Waiting for Docker services to start (15 seconds)...$(RESET)"
	@sleep 15
	@echo "$(GREEN)Starting server in background...$(RESET)"
	@$(GO) run $(SERVER_CMD) > server.log 2>&1 &
	@echo "$(YELLOW)Waiting for server to initialize (5 seconds)...$(RESET)"
	@sleep 5
	@echo "$(GREEN)Starting stock bot in background...$(RESET)"
	@$(GO) run $(BOT_CMD) > bot.log 2>&1 &
	@echo "$(GREEN)All services restarted with fresh database!$(RESET)"
	@echo "$(GREEN)Access the application at: http://localhost:8080$(RESET)"
	@echo "$(YELLOW)View server logs: tail -f server.log$(RESET)"
	@echo "$(YELLOW)View bot logs: tail -f bot.log$(RESET)"
	@echo "$(GREEN)Register a new user to begin testing!$(RESET)"

# Stop all services
stop-all:
	@echo "$(YELLOW)Stopping Go processes...$(RESET)"
	@pkill -f "go run" || true
	@echo "$(YELLOW)Stopping Docker containers...$(RESET)"
	@$(DOCKER) down
	@echo "$(GREEN)All services stopped!$(RESET)"

# Run server and bot using Docker for dependencies and local execution for the app
run:
	@echo "$(GREEN)Starting PostgreSQL and RabbitMQ with Docker...$(RESET)"
	@$(DOCKER) up -d postgres rabbitmq
	@echo "$(YELLOW)Waiting for services to start...$(RESET)"
	@sleep 10
	@echo "$(GREEN)Starting local server on port 8080...$(RESET)"
	@$(GO) run $(SERVER_CMD) & \
	echo "$(GREEN)Starting local stock bot...$(RESET)" && \
	$(GO) run $(BOT_CMD)

# Run local version (legacy)
run-local:
	@echo "$(GREEN)Starting application locally...$(RESET)"
	@echo "$(YELLOW)Make sure PostgreSQL and RabbitMQ are running!$(RESET)"
	@echo "$(GREEN)Starting server on port 8080...$(RESET)"
	@$(GO) run $(SERVER_CMD) & \
	echo "$(GREEN)Starting stock bot...$(RESET)" && \
	$(GO) run $(BOT_CMD)

# Run only the web server
run-server:
	@echo "$(GREEN)Starting web server on port 8080...$(RESET)"
	@$(GO) run $(SERVER_CMD)

# Run only the stock bot
run-bot:
	@echo "$(GREEN)Starting stock bot...$(RESET)"
	@$(GO) run $(BOT_CMD)

# Build binaries
build:
	@echo "$(GREEN)Compiling application...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/server $(SERVER_CMD)
	@$(GO) build -o $(BUILD_DIR)/bot $(BOT_CMD)
	@echo "$(GREEN)Binaries created in $(BUILD_DIR)/$(RESET)"

# Run tests
test:
	@echo "$(GREEN)Running tests...$(RESET)"
	@$(GO) test ./... -v

# Clean binaries and temporary files
clean:
	@echo "$(YELLOW)Removing binaries and temporary files...$(RESET)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)Cleanup completed.$(RESET)"

# Help
help:
	@echo "$(GREEN)Available commands:$(RESET)"
	@echo "  make run-all      - Start all services (main command to run the application)"
	@echo "  make stop-all     - Stop all services (Docker and Go processes)"
	@echo "  make reset-all    - Reset and restart all services with fresh database (use when you need a clean start)"
	@echo "  make run          - Start all services with Docker Compose (foreground with logs)"
	@echo "  make docker       - Start all services with Docker Compose (detached mode)"
	@echo "  make docker-down  - Stop all Docker containers"
	@echo "  make run-local    - Run server and bot locally (requires PostgreSQL and RabbitMQ)"
	@echo "  make run-server   - Run only the web server locally"
	@echo "  make run-bot      - Run only the stock bot locally"
	@echo "  make build        - Compile application binaries"
	@echo "  make test         - Run all tests"
	@echo "  make clean        - Remove binaries and temporary files" 