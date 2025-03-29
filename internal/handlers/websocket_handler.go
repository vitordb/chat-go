package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/dbvitor/chat-go/internal/models"
	"github.com/dbvitor/chat-go/internal/services"
	"github.com/dbvitor/chat-go/pkg/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

// WebSocketHandler handles WebSocket connections for real-time chat
type WebSocketHandler struct {
	messageService  *services.MessageService
	userService     *services.UserService
	chatroomService *services.ChatroomService
	clients         map[string]map[*websocket.Conn]bool // Map of chatroom ID to client connections
	clientsMutex    sync.RWMutex
	upgrader        websocket.Upgrader
	stockResults    <-chan amqp.Delivery
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(
	messageService *services.MessageService,
	userService *services.UserService,
	chatroomService *services.ChatroomService,
	stockResults <-chan amqp.Delivery,
) *WebSocketHandler {
	handler := &WebSocketHandler{
		messageService:  messageService,
		userService:     userService,
		chatroomService: chatroomService,
		clients:         make(map[string]map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for this example
			},
		},
		stockResults: stockResults,
	}

	// Start processing stock results
	go handler.processStockResults()

	return handler
}

// MessagePayload represents a message sent over the WebSocket connection
type MessagePayload struct {
	Content string `json:"content"`
}

// Handle upgrades HTTP connection to WebSocket and manages communication
func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated
	if !auth.IsAuthenticated(r) {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Get user ID
	userID, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Get chatroom ID from URL
	vars := mux.Vars(r)
	chatroomID := vars["id"]

	// Check if chatroom exists
	_, err = h.chatroomService.GetByID(chatroomID)
	if err != nil {
		http.Error(w, "Chatroom not found", http.StatusNotFound)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Register client
	h.clientsMutex.Lock()
	if _, ok := h.clients[chatroomID]; !ok {
		h.clients[chatroomID] = make(map[*websocket.Conn]bool)
	}
	h.clients[chatroomID][conn] = true
	h.clientsMutex.Unlock()

	// Get messages for chatroom
	messages, err := h.messageService.GetMessagesByChatroomID(chatroomID)
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
	} else {
		// Send messages to client
		for _, message := range messages {
			if err := conn.WriteJSON(message); err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}

	// Send system message
	user, _ := h.userService.GetByID(userID)
	systemMessage := models.NewMessage("", "System", chatroomID, user.Username+" joined the chat", models.MessageTypeSystem)
	h.broadcastMessage(systemMessage, chatroomID)

	// Handle incoming messages
	go h.handleClient(conn, userID, chatroomID)
}

// handleClient processes messages from a WebSocket client
func (h *WebSocketHandler) handleClient(conn *websocket.Conn, userID, chatroomID string) {
	defer func() {
		// Unregister client
		h.clientsMutex.Lock()
		delete(h.clients[chatroomID], conn)
		h.clientsMutex.Unlock()
		conn.Close()

		// Send system message
		user, err := h.userService.GetByID(userID)
		if err == nil {
			systemMessage := models.NewMessage("", "System", chatroomID, user.Username+" left the chat", models.MessageTypeSystem)
			h.broadcastMessage(systemMessage, chatroomID)
		}
	}()

	for {
		// Read message from client
		var payload MessagePayload
		err := conn.ReadJSON(&payload)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Create and save message
		message, err := h.messageService.CreateMessage(userID, chatroomID, payload.Content)
		if err != nil {
			log.Printf("Error creating message: %v", err)
			continue
		}

		// If message is nil, it's a stock command and doesn't need to be broadcast
		if message != nil {
			// Broadcast message to all clients in the chatroom
			h.broadcastMessage(message, chatroomID)
		}
	}
}

// broadcastMessage sends a message to all clients in a chatroom
func (h *WebSocketHandler) broadcastMessage(message *models.Message, chatroomID string) {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()

	for client := range h.clients[chatroomID] {
		err := client.WriteJSON(message)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
			client.Close()
			delete(h.clients[chatroomID], client)
		}
	}
}

// processStockResults listens for stock results and broadcasts them
func (h *WebSocketHandler) processStockResults() {
	for delivery := range h.stockResults {
		// Parse message
		var result map[string]interface{}
		err := json.Unmarshal(delivery.Body, &result)
		if err != nil {
			log.Printf("Error parsing stock result: %v", err)
			continue
		}

		// Extract data
		chatroomID, ok := result["chatroom_id"].(string)
		if !ok {
			log.Println("Invalid chatroom ID in stock result")
			continue
		}

		// Create stock response
		stockResponse := &models.StockResponse{}

		// Extract symbol
		if symbol, ok := result["symbol"].(string); ok {
			stockResponse.Symbol = symbol
		}

		// Extract price
		if price, ok := result["price"].(float64); ok {
			stockResponse.Price = price
		}

		// Extract error
		if errMsg, ok := result["error"].(string); ok {
			stockResponse.Error = errMsg
		}

		// Create bot message
		message, err := h.messageService.CreateBotMessage(chatroomID, stockResponse)
		if err != nil {
			log.Printf("Error creating bot message: %v", err)
			continue
		}

		// Broadcast message
		h.broadcastMessage(message, chatroomID)
	}
}
