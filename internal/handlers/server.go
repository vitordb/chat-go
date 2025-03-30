package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dbvitor/chat-go/internal/services"
	"github.com/dbvitor/chat-go/pkg/broker"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Server represents the HTTP server
type Server struct {
	router          *mux.Router
	userHandler     *UserHandler
	chatroomHandler *ChatroomHandler
	wsHandler       *WebSocketHandler
}

// NewServer creates a new HTTP server
func NewServer(db *sql.DB, rabbitMQ *broker.RabbitMQ, stockResults <-chan amqp.Delivery) *Server {
	// Create services
	userService := services.NewUserService(db)
	chatroomService := services.NewChatroomService(db)
	messageService := services.NewMessageService(db, rabbitMQ)

	// Create handlers
	userHandler := NewUserHandler(userService)
	chatroomHandler := NewChatroomHandler(chatroomService)
	wsHandler := NewWebSocketHandler(messageService, userService, chatroomService, stockResults)

	// Create router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(corsMiddleware)

	// Register API routes
	apiRouter := router.PathPrefix("/api").Subrouter()

	// User routes
	apiRouter.HandleFunc("/auth/register", userHandler.Register).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/auth/login", userHandler.Login).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/auth/logout", userHandler.Logout).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/auth/check", userHandler.CheckAuth).Methods("GET", "OPTIONS")

	// Chatroom routes
	apiRouter.HandleFunc("/chatrooms", chatroomHandler.GetAll).Methods("GET", "OPTIONS")
	apiRouter.HandleFunc("/chatrooms", chatroomHandler.Create).Methods("POST", "OPTIONS")
	apiRouter.HandleFunc("/chatrooms/{id}", chatroomHandler.GetByID).Methods("GET", "OPTIONS")

	// WebSocket route
	apiRouter.HandleFunc("/ws/{id}", wsHandler.Handle)

	// Static files
	fs := http.FileServer(http.Dir("./web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Serve index.html for all other routes
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "index.html"))
	})

	return &Server{
		router:          router,
		userHandler:     userHandler,
		chatroomHandler: chatroomHandler,
		wsHandler:       wsHandler,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	return http.ListenAndServe(":"+port, s.router)
}
