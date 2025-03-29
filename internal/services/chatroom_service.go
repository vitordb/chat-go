package services

import (
	"database/sql"

	"github.com/dbvitor/chat-go/internal/database"
	"github.com/dbvitor/chat-go/internal/models"
)

// ChatroomService handles chatroom-related business logic
type ChatroomService struct {
	chatroomRepo *database.ChatroomRepository
}

// NewChatroomService creates a new chatroom service
func NewChatroomService(db *sql.DB) *ChatroomService {
	return &ChatroomService{
		chatroomRepo: database.NewChatroomRepository(db),
	}
}

// Create creates a new chatroom
func (s *ChatroomService) Create(name string) (*models.Chatroom, error) {
	// Create new chatroom
	chatroom := models.NewChatroom(name)

	// Save chatroom to database
	err := s.chatroomRepo.Create(chatroom)
	if err != nil {
		return nil, err
	}

	return chatroom, nil
}

// GetByID retrieves a chatroom by ID
func (s *ChatroomService) GetByID(id string) (*models.Chatroom, error) {
	return s.chatroomRepo.GetByID(id)
}

// GetByName retrieves a chatroom by name
func (s *ChatroomService) GetByName(name string) (*models.Chatroom, error) {
	return s.chatroomRepo.GetByName(name)
}

// GetAll retrieves all chatrooms
func (s *ChatroomService) GetAll() ([]*models.Chatroom, error) {
	return s.chatroomRepo.GetAll()
}
