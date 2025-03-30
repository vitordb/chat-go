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

// Maximum number of messages to load from history
const MaxMessages = 50

// Regex to validate stock quote commands
var StockCommandPattern = regexp.MustCompile(`^/stock=([A-Za-z0-9.]+)$`)

// Service responsible for message operations
type MessageService struct {
	messageRepo *database.MessageRepository
	userRepo    *database.UserRepository
	rabbitMQ    *broker.RabbitMQ
}

// Creates a new instance of the message service
func NewMessageService(db *sql.DB, rabbitMQ *broker.RabbitMQ) *MessageService {
	return &MessageService{
		messageRepo: database.NewMessageRepository(db),
		userRepo:    database.NewUserRepository(db),
		rabbitMQ:    rabbitMQ,
	}
}

// Creates a new message and saves it to the database, or processes special commands
func (s *MessageService) CreateMessage(userID, chatroomID, content string) (*models.Message, error) {
	// Get user data
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if it's a stock command
	if match := StockCommandPattern.FindStringSubmatch(content); match != nil {
		stockCode := match[1]

		// Send request to the stock bot
		err := s.rabbitMQ.PublishStockRequest(chatroomID, stockCode)
		if err != nil {
			return nil, err
		}

		// We don't save the command as a message
		return nil, nil
	}

	// Create normal message
	message := models.NewMessage(userID, user.Username, chatroomID, content, models.MessageTypeChat)

	// Save to database
	err = s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Creates a message from the stock bot
func (s *MessageService) CreateBotMessage(chatroomID string, stockResponse *models.StockResponse) (*models.Message, error) {
	var content string

	if stockResponse.Error != "" {
		// Error message for quote
		content = fmt.Sprintf("Error getting quote for %s: %s", stockResponse.Symbol, stockResponse.Error)
	} else {
		// Message with the quote
		content = fmt.Sprintf("%s quote is $%.2f per share", strings.ToUpper(stockResponse.Symbol), stockResponse.Price)
	}

	// Create bot message using the bot ID instead of empty string
	message := models.NewMessage(database.BotUserID, "Stock Bot", chatroomID, content, models.MessageTypeStock)

	// Save to database
	err := s.messageRepo.Create(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Gets messages for a specific chatroom
func (s *MessageService) GetMessagesByChatroomID(chatroomID string) ([]*models.Message, error) {
	return s.messageRepo.GetByChatroomID(chatroomID, MaxMessages)
}
