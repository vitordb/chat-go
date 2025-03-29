package models

import (
	"time"
)

// Chatroom represents a chat room where users can send messages
type Chatroom struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewChatroom creates a new chat room
func NewChatroom(name string) *Chatroom {
	now := time.Now()
	return &Chatroom{
		ID:        "",
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
