package models

import (
	"time"
)

type MessageType string

const (
	MessageTypeChat   MessageType = "chat"
	MessageTypeStock  MessageType = "stock"
	MessageTypeSystem MessageType = "system"
)

type Message struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Username   string      `json:"username"`
	ChatroomID string      `json:"chatroom_id"`
	Content    string      `json:"content"`
	Type       MessageType `json:"type"`
	CreatedAt  time.Time   `json:"created_at"`
}

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

type StockResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Error  string  `json:"error,omitempty"`
}
