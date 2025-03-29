package auth

import (
	"errors"
	"net/http"
	"os"

	"github.com/dbvitor/chat-go/internal/models"
	"github.com/gorilla/sessions"
)

// Define error types
var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrNotAuthenticated   = errors.New("not authenticated")
)

// SessionName is the name of the cookie used to store the session
const SessionName = "chat-session"

// UserKey is the key used to store the user in the session
const UserKey = "user_id"

// Store is the session store
var Store *sessions.CookieStore

// Initialize sets up the authentication module
func Initialize() {
	// Create a new cookie store
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	// Set session options
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 1 week
		HttpOnly: true,
	}
}

// Authenticate verifies user credentials and creates a session
func Authenticate(w http.ResponseWriter, r *http.Request, user *models.User) error {
	// Create a new session
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	// Store user ID in session
	session.Values[UserKey] = user.ID

	// Save the session
	return session.Save(r, w)
}

// GetAuthenticatedUser retrieves the authenticated user from the session
func GetAuthenticatedUser(r *http.Request) (string, error) {
	// Get the session
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return "", err
	}

	// Get the user ID from the session
	userID, ok := session.Values[UserKey].(string)
	if !ok {
		return "", ErrNotAuthenticated
	}

	return userID, nil
}

// Logout removes the user from the session
func Logout(w http.ResponseWriter, r *http.Request) error {
	// Get the session
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return err
	}

	// Delete the user ID from the session
	delete(session.Values, UserKey)

	// Set session to expire
	session.Options.MaxAge = -1

	// Save the session
	return session.Save(r, w)
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	// Get the session
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return false
	}

	// Check if user ID exists in session
	_, ok := session.Values[UserKey].(string)
	return ok
}
