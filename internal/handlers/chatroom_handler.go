package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dbvitor/chat-go/internal/services"
	"github.com/dbvitor/chat-go/pkg/auth"
	"github.com/gorilla/mux"
)

// ChatroomHandler handles chatroom-related HTTP requests
type ChatroomHandler struct {
	chatroomService *services.ChatroomService
}

// NewChatroomHandler creates a new chatroom handler
func NewChatroomHandler(chatroomService *services.ChatroomService) *ChatroomHandler {
	return &ChatroomHandler{
		chatroomService: chatroomService,
	}
}

// CreateChatroomRequest represents the request body for creating a chatroom
type CreateChatroomRequest struct {
	Name string `json:"name"`
}

// Create handles chatroom creation
func (h *ChatroomHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated
	if !auth.IsAuthenticated(r) {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req CreateChatroomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" {
		http.Error(w, "Chatroom name is required", http.StatusBadRequest)
		return
	}

	// Create chatroom
	chatroom, err := h.chatroomService.Create(req.Name)
	if err != nil {
		http.Error(w, "Failed to create chatroom", http.StatusInternalServerError)
		return
	}

	// Return chatroom info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chatroom)
}

// GetAll handles retrieving all chatrooms
func (h *ChatroomHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated
	if !auth.IsAuthenticated(r) {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Get all chatrooms
	chatrooms, err := h.chatroomService.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve chatrooms", http.StatusInternalServerError)
		return
	}

	// Return chatrooms
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatrooms)
}

// GetByID handles retrieving a chatroom by ID
func (h *ChatroomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated
	if !auth.IsAuthenticated(r) {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Get chatroom ID from URL
	vars := mux.Vars(r)
	chatroomID := vars["id"]

	// Get chatroom
	chatroom, err := h.chatroomService.GetByID(chatroomID)
	if err != nil {
		http.Error(w, "Chatroom not found", http.StatusNotFound)
		return
	}

	// Return chatroom info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatroom)
}
