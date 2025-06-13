// Package api contains the main API server implementation
package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Asif-Faizal/Gommerce/services/user"
	"github.com/gorilla/mux" // Popular HTTP router for Go
)

// APIServer represents our main server structure
// It holds the server configuration and database connection
type APIServer struct {
	listenAddress string  // The address where the server will listen (e.g., ":3000")
	db            *sql.DB // Database connection pointer
}

// NewAPIServer creates a new instance of APIServer
// It's a constructor function that initializes the server with given parameters
func NewAPIServer(listenAddress string, db *sql.DB) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		db:            db,
	}
}

// Run starts the HTTP server and sets up all routes
// Returns an error if the server fails to start
func (s *APIServer) Run() error {
	// Create a new router instance
	router := mux.NewRouter()

	// Create a subrouter for API versioning
	// All routes will be prefixed with /api/v1
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	log.Println("Starting server on", s.listenAddress)

	// Initialize user handler and register its routes
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Start the HTTP server and listen for incoming requests
	return http.ListenAndServe(s.listenAddress, router)
}
