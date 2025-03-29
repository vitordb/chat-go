# Makefile para o projeto Chat-Go

# Variáveis
GO=go
DOCKER_COMPOSE=docker-compose
SERVER_CMD=./cmd/server/main.go
BOT_CMD=./cmd/bot/main.go
BUILD_DIR=./build

# Cores para melhorar a legibilidade
VERDE=\033[0;32m
AMARELO=\033[1;33m
RESET=\033[0m

# Comandos principais
.PHONY: all run run-server run-bot build clean test docker docker-down help

all: build

# Iniciar toda a aplicação com Docker Compose
docker:
	@echo "$(VERDE)Iniciando aplicação com Docker Compose...$(RESET)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(VERDE)Serviços iniciados! Acesse http://localhost:8080$(RESET)"

# Parar os containers Docker
docker-down:
	@echo "$(AMARELO)Parando todos os containers...$(RESET)"
	@$(DOCKER_COMPOSE) down
	@echo "$(VERDE)Containers parados com sucesso.$(RESET)"

# Executar o servidor e o bot em terminais separados
run:
	@echo "$(VERDE)Iniciando aplicação...$(RESET)"
	@echo "$(AMARELO)Certifique-se de que PostgreSQL e RabbitMQ estejam funcionando!$(RESET)"
	@echo "$(VERDE)Iniciando servidor na porta 8080...$(RESET)"
	@$(GO) run $(SERVER_CMD) & \
	echo "$(VERDE)Iniciando bot de cotações...$(RESET)" && \
	$(GO) run $(BOT_CMD)

# Executar apenas o servidor web
run-server:
	@echo "$(VERDE)Iniciando servidor web na porta 8080...$(RESET)"
	@$(GO) run $(SERVER_CMD)

# Executar apenas o bot de cotações
run-bot:
	@echo "$(VERDE)Iniciando bot de cotações...$(RESET)"
	@$(GO) run $(BOT_CMD)

# Construir os binários
build:
	@echo "$(VERDE)Compilando aplicação...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/server $(SERVER_CMD)
	@$(GO) build -o $(BUILD_DIR)/bot $(BOT_CMD)
	@echo "$(VERDE)Binários criados em $(BUILD_DIR)/$(RESET)"

# Executar os testes
test:
	@echo "$(VERDE)Executando testes...$(RESET)"
	@$(GO) test ./... -v

# Limpar binários e arquivos temporários
clean:
	@echo "$(AMARELO)Removendo binários e arquivos temporários...$(RESET)"
	@rm -rf $(BUILD_DIR)
	@echo "$(VERDE)Limpeza concluída.$(RESET)"

# Ajuda
help:
	@echo "$(VERDE)Comandos disponíveis:$(RESET)"
	@echo "  make docker       - Inicia todos os serviços com Docker Compose"
	@echo "  make docker-down  - Para todos os containers Docker"
	@echo "  make run          - Executa o servidor e o bot simultaneamente"
	@echo "  make run-server   - Executa apenas o servidor web"
	@echo "  make run-bot      - Executa apenas o bot de cotações"
	@echo "  make build        - Compila os binários da aplicação"
	@echo "  make test         - Executa todos os testes"
	@echo "  make clean        - Remove binários e arquivos temporários" 