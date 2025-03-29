package services

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/dbvitor/chat-go/internal/database"
	"github.com/dbvitor/chat-go/internal/models"
	"github.com/dbvitor/chat-go/pkg/broker"
)

// Max messages to retrieve from the database
const MaxMessages = 50

// StockCommandPattern is the pattern for stock commands
var StockCommandPattern = regexp.MustCompile(`^/stock=([A-Za-z0-9.]+)$`)

// MessageService handles message-related business logic
type MessageService struct {
	messageRepo *database.MessageRepository
	userRepo    *database.UserRepository
	rabbitMQ    *broker.RabbitMQ
}

// NewMessageService creates a new message service
func NewMessageService(db *sql.DB, rabbitMQ *broker.RabbitMQ) *MessageService {
	return &MessageService{
		messageRepo: database.NewMessageRepository(db),
		userRepo:    database.NewUserRepository(db),
		rabbitMQ:    rabbitMQ,
	}
}

// CreateMessage creates a new message
func (s *MessageService) CreateMessage(userID, chatroomID, content string) (*models.Message, error) {
	// Get user by ID
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if message is a stock command
	if match := StockCommandPattern.FindStringSubmatch(content); match != nil {
		// Extract stock code
		stockCode := match[1]

		// Publish stock request
		err := s.rabbitMQ.PublishStockRequest(chatroomID, stockCode)
		if err != nil {
			return nil, err
		}

		// Return nil to indicate that the message should not be stored
		return nil, nil
	}

	// Create new message
	message := models.NewMessage(userID, user.Username, chatroomID, content, models.MessageTypeChat)

	// Save message to database
	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// CreateBotMessage creates a message from the stock bot
func (s *MessageService) CreateBotMessage(chatroomID string, stockResponse *models.StockResponse) (*models.Message, error) {
	var content string

	if stockResponse.Error != "" {
		// Create error message
		content = fmt.Sprintf("Error getting quote for %s: %s", stockResponse.Symbol, stockResponse.Error)
	} else {
		// Create stock quote message
		content = fmt.Sprintf("%s quote is $%.2f per share", strings.ToUpper(stockResponse.Symbol), stockResponse.Price)
	}

	// Create new message
	message := models.NewMessage("", "Stock Bot", chatroomID, content, models.MessageTypeStock)

	// Save message to database
	err := s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessagesByChatroomID retrieves messages for a specific chatroom
func (s *MessageService) GetMessagesByChatroomID(chatroomID string) ([]*models.Message, error) {
	return s.messageRepo.GetByChatroomID(chatroomID, MaxMessages)
}
