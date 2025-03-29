package services

import (
	"database/sql"
	"errors"

	"github.com/dbvitor/chat-go/internal/database"
	"github.com/dbvitor/chat-go/internal/models"
	"github.com/dbvitor/chat-go/pkg/auth"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo *database.UserRepository
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		userRepo: database.NewUserRepository(db),
	}
}

// Register creates a new user
func (s *UserService) Register(username, password string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existingUser != nil {
		return nil, auth.ErrUserAlreadyExists
	}

	// Create new user
	user, err := models.NewUser(username, password)
	if err != nil {
		return nil, err
	}

	// Save user to database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login validates user credentials
func (s *UserService) Login(username, password string) (*models.User, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, auth.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, auth.ErrInvalidCredentials
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id string) (*models.User, error) {
	return s.userRepo.GetByID(id)
}
