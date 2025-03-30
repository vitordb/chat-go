package database

import (
	"database/sql"
	"time"

	"github.com/dbvitor/chat-go/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, password, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return err
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password, created_at, updated_at 
	          FROM users 
	          WHERE username = $1`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `SELECT id, username, password, created_at, updated_at 
	          FROM users 
	          WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users 
	          SET username = $1, password = $2, updated_at = $3 
	          WHERE id = $4`

	_, err := r.db.Exec(
		query,
		user.Username,
		user.Password,
		time.Now(),
		user.ID,
	)

	return err
}

func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
