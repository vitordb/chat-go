package models

import (
	"time"
)

// MessageType defines the type of message
type MessageType string

const (
	// MessageTypeChat is a regular chat message
	MessageTypeChat MessageType = "chat"
	// MessageTypeStock is a stock quote message
	MessageTypeStock MessageType = "stock"
	// MessageTypeSystem is a system notification
	MessageTypeSystem MessageType = "system"
)

// Message represents a chat message
type Message struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Username   string      `json:"username"`
	ChatroomID string      `json:"chatroom_id"`
	Content    string      `json:"content"`
	Type       MessageType `json:"type"`
	CreatedAt  time.Time   `json:"created_at"`
}

// NewMessage creates a new chat message
func NewMessage(userID, username, chatroomID, content string, msgType MessageType) *Message {
	return &Message{
		ID:         "",
		UserID:     userID,
		Username:   username,
		ChatroomID: chatroomID,
		Content:    content,
		Type:       msgType,
		CreatedAt:  time.Now(),
	}
}

// StockResponse represents a response from the stock API
type StockResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Error  string  `json:"error,omitempty"`
}
