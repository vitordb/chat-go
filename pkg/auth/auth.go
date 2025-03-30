package auth

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/dbvitor/chat-go/internal/models"
	"github.com/gorilla/sessions"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrNotAuthenticated   = errors.New("not authenticated")
)

const SessionName = "chat-session"

const UserKey = "user_id"

var Store *sessions.CookieStore

// Initialize sets up the authentication module
func Initialize() {
	// Get session key from environment variable, or use a default key
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "supersecretkey123456789"
		log.Println("Warning: Using default session key. Set SESSION_KEY for better security.")
	}

	// Create a new cookie store
	Store = sessions.NewCookieStore([]byte(sessionKey))

	// Set session options
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		// Add SameSite attribute for better security
		SameSite: http.SameSiteLaxMode,
	}
}

// Authenticate verifies user credentials and creates a session
func Authenticate(w http.ResponseWriter, r *http.Request, user *models.User) error {
	// Create a new session
	session, err := Store.Get(r, SessionName)
	if err != nil {
		log.Printf("Session error: %v", err)
		session = sessions.NewSession(Store, SessionName)
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
	}

	// Store user ID in session
	session.Values[UserKey] = user.ID

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session: %v", err)
	}
	return err
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
	session, err := Store.Get(r, SessionName)
	if err != nil {
		return false
	}

	// Check if user ID exists in session
	_, ok := session.Values[UserKey].(string)
	return ok
}
