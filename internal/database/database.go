package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// DB is the database connection
var DB *sql.DB

// BotUserID is the ID of the bot user in the database
const BotUserID = "00000000-0000-0000-0000-000000000000"

// Initialize sets up the database connection
func Initialize() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	return createTables()
}

// createTables creates the necessary tables if they don't exist
func createTables() error {
	// Create users table
	userTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err := DB.Exec(userTableQuery)
	if err != nil {
		return err
	}

	// Create chatrooms table
	chatroomTableQuery := `
	CREATE TABLE IF NOT EXISTS chatrooms (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(50) UNIQUE NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err = DB.Exec(chatroomTableQuery)
	if err != nil {
		return err
	}

	// Create messages table
	messageTableQuery := `
	CREATE TABLE IF NOT EXISTS messages (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id),
		username VARCHAR(50) NOT NULL,
		chatroom_id UUID NOT NULL REFERENCES chatrooms(id),
		content TEXT NOT NULL,
		type VARCHAR(20) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`

	_, err = DB.Exec(messageTableQuery)
	if err != nil {
		return err
	}

	// Insert default chatroom if none exists
	_, err = DB.Exec("INSERT INTO chatrooms (name) VALUES ('General') ON CONFLICT DO NOTHING;")
	if err != nil {
		return err
	}

	// Insert the bot user if it doesn't exist
	botUserQuery := `
	INSERT INTO users (id, username, password, created_at, updated_at)
	VALUES ($1, 'Stock Bot', 'botpassword', NOW(), NOW())
	ON CONFLICT (id) DO NOTHING;
	`
	_, err = DB.Exec(botUserQuery, BotUserID)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
