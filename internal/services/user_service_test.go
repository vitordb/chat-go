package services

import (
	"errors"
	"testing"

	"github.com/dbvitor/chat-go/internal/models"
	"github.com/dbvitor/chat-go/pkg/auth"
)

// MockUserRepository for testing
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *MockUserRepository) Create(user *models.User) error {
	// Check if user already exists
	for _, u := range r.users {
		if u.Username == user.Username {
			return errors.New("user already exists")
		}
	}

	user.ID = "test-id"
	r.users[user.Username] = user
	return nil
}

func (r *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *MockUserRepository) GetByID(id string) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// UserService with mock repository
type TestUserService struct {
	repo *MockUserRepository
}

func NewTestUserService() *TestUserService {
	return &TestUserService{
		repo: NewMockUserRepository(),
	}
}

func (s *TestUserService) Register(username, password string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetByUsername(username)
	if err == nil && existingUser != nil {
		return nil, auth.ErrUserAlreadyExists
	}

	// Create new user
	user, err := models.NewUser(username, password)
	if err != nil {
		return nil, err
	}

	// Save user to mock repository
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *TestUserService) Login(username, password string) (*models.User, error) {
	// Get user by username
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, auth.ErrInvalidCredentials
	}

	return user, nil
}

func (s *TestUserService) GetByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

// Tests
func TestUserService_Register(t *testing.T) {
	service := NewTestUserService()

	// Test registration with new user
	user, err := service.Register("testuser", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got: %s", user.Username)
	}

	// Test registration with existing user
	_, err = service.Register("testuser", "password123")
	if err != auth.ErrUserAlreadyExists {
		t.Errorf("Expected error 'user already exists', got: %v", err)
	}
}

func TestUserService_Login(t *testing.T) {
	service := NewTestUserService()

	// Register a test user first
	testUser, _ := models.NewUser("testuser", "password123")
	testUser.ID = "test-id"
	service.repo.users[testUser.Username] = testUser

	// Test login with correct credentials
	user, err := service.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got: %s", user.Username)
	}

	// Test login with incorrect password
	_, err = service.Login("testuser", "wrongpassword")
	if err != auth.ErrInvalidCredentials {
		t.Errorf("Expected error 'invalid credentials', got: %v", err)
	}

	// Test login with non-existent user
	_, err = service.Login("nonexistent", "password123")
	if err != auth.ErrInvalidCredentials {
		t.Errorf("Expected error 'invalid credentials', got: %v", err)
	}
}
