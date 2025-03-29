package database

import (
	"database/sql"

	"github.com/dbvitor/chat-go/internal/models"
)

// MessageRepository handles message database operations
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create adds a new message to the database
func (r *MessageRepository) Create(message *models.Message) error {
	query := `INSERT INTO messages (user_id, username, chatroom_id, content, type, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6) 
	          RETURNING id`

	err := r.db.QueryRow(
		query,
		message.UserID,
		message.Username,
		message.ChatroomID,
		message.Content,
		message.Type,
		message.CreatedAt,
	).Scan(&message.ID)

	return err
}

// GetByChatroomID retrieves messages for a specific chatroom, limited to the last 50
func (r *MessageRepository) GetByChatroomID(chatroomID string, limit int) ([]*models.Message, error) {
	query := `SELECT id, user_id, username, chatroom_id, content, type, created_at 
	          FROM messages 
	          WHERE chatroom_id = $1 
	          ORDER BY created_at DESC 
	          LIMIT $2`

	rows, err := r.db.Query(query, chatroomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Username,
			&message.ChatroomID,
			&message.Content,
			&message.Type,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Reverse the messages to get chronological order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetByID retrieves a message by ID
func (r *MessageRepository) GetByID(id string) (*models.Message, error) {
	query := `SELECT id, user_id, username, chatroom_id, content, type, created_at 
	          FROM messages 
	          WHERE id = $1`

	var message models.Message
	err := r.db.QueryRow(query, id).Scan(
		&message.ID,
		&message.UserID,
		&message.Username,
		&message.ChatroomID,
		&message.Content,
		&message.Type,
		&message.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// Delete removes a message from the database
func (r *MessageRepository) Delete(id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
