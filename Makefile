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
RESET=\033[0m

# Main commands
.PHONY: all run run-server run-bot build clean test docker docker-down help

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

# Run server and bot in separate terminals
run:
	@echo "$(GREEN)Starting application...$(RESET)"
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
	@echo "  make docker       - Start all services with Docker Compose"
	@echo "  make docker-down  - Stop all Docker containers"
	@echo "  make run          - Run server and bot simultaneously"
	@echo "  make run-server   - Run only the web server"
	@echo "  make run-bot      - Run only the stock bot"
	@echo "  make build        - Compile application binaries"
	@echo "  make test         - Run all tests"
	@echo "  make clean        - Remove binaries and temporary files" 