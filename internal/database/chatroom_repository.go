package database

import (
	"database/sql"
	"time"

	"github.com/dbvitor/chat-go/internal/models"
)

// ChatroomRepository handles chatroom database operations
type ChatroomRepository struct {
	db *sql.DB
}

// NewChatroomRepository creates a new chatroom repository
func NewChatroomRepository(db *sql.DB) *ChatroomRepository {
	return &ChatroomRepository{db: db}
}

// Create adds a new chatroom to the database
func (r *ChatroomRepository) Create(chatroom *models.Chatroom) error {
	query := `INSERT INTO chatrooms (name, created_at, updated_at) 
	          VALUES ($1, $2, $3) 
	          RETURNING id`

	err := r.db.QueryRow(
		query,
		chatroom.Name,
		chatroom.CreatedAt,
		chatroom.UpdatedAt,
	).Scan(&chatroom.ID)

	return err
}

// GetByID retrieves a chatroom by ID
func (r *ChatroomRepository) GetByID(id string) (*models.Chatroom, error) {
	query := `SELECT id, name, created_at, updated_at 
	          FROM chatrooms 
	          WHERE id = $1`

	var chatroom models.Chatroom
	err := r.db.QueryRow(query, id).Scan(
		&chatroom.ID,
		&chatroom.Name,
		&chatroom.CreatedAt,
		&chatroom.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &chatroom, nil
}

// GetByName retrieves a chatroom by name
func (r *ChatroomRepository) GetByName(name string) (*models.Chatroom, error) {
	query := `SELECT id, name, created_at, updated_at 
	          FROM chatrooms 
	          WHERE name = $1`

	var chatroom models.Chatroom
	err := r.db.QueryRow(query, name).Scan(
		&chatroom.ID,
		&chatroom.Name,
		&chatroom.CreatedAt,
		&chatroom.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &chatroom, nil
}

// GetAll retrieves all chatrooms
func (r *ChatroomRepository) GetAll() ([]*models.Chatroom, error) {
	query := `SELECT id, name, created_at, updated_at 
	          FROM chatrooms 
	          ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatrooms []*models.Chatroom
	for rows.Next() {
		var chatroom models.Chatroom
		err := rows.Scan(
			&chatroom.ID,
			&chatroom.Name,
			&chatroom.CreatedAt,
			&chatroom.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		chatrooms = append(chatrooms, &chatroom)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chatrooms, nil
}

// Update updates chatroom information
func (r *ChatroomRepository) Update(chatroom *models.Chatroom) error {
	query := `UPDATE chatrooms 
	          SET name = $1, updated_at = $2 
	          WHERE id = $3`

	_, err := r.db.Exec(
		query,
		chatroom.Name,
		time.Now(),
		chatroom.ID,
	)

	return err
}

// Delete removes a chatroom from the database
func (r *ChatroomRepository) Delete(id string) error {
	query := `DELETE FROM chatrooms WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
