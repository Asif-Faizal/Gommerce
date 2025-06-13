// Package user contains the user-related HTTP handlers and routes
package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Asif-Faizal/Gommerce/services/auth"
	"github.com/Asif-Faizal/Gommerce/types"
	"github.com/Asif-Faizal/Gommerce/utils"
	"github.com/gorilla/mux"
)

// Handler represents the user-related HTTP handlers
// It contains methods to handle different user-related endpoints
type Handler struct {
	store types.UserStore // Interface for user data operations
}

// NewHandler creates a new instance of the user Handler
// This is a constructor function for the Handler struct
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes sets up all the user-related routes
// It takes a router and attaches the handler functions to specific paths
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Register the login endpoint - will handle POST requests to /api/v1/login
	router.HandleFunc("/login", h.handleLogin)

	// Register the registration endpoint - will handle POST requests to /api/v1/register
	router.HandleFunc("/register", h.handleRegister)
}

// handleLogin processes user login requests
// w is the response writer to send back HTTP responses
// r is the HTTP request containing the login data
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement login logic
}

// handleRegister processes user registration requests
// w is the response writer to send back HTTP responses
// r is the HTTP request containing the registration data
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validate the payload
	if err := h.validateRegisterPayload(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Check if user already exists
	existingUser, err := h.store.GetUserByEmail(payload.Email)
	if err != nil && err != sql.ErrNoRows {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error checking user existence: %w", err))
		return
	}
	if existingUser != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error hashing password: %w", err))
		return
	}

	// Create new user
	user := &types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	// Save user to database
	if err := h.store.CreateUser(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating user: %w", err))
		return
	}

	// Return success response
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "user created successfully",
	})
}

// validateRegisterPayload validates the registration payload
// Returns an error if any required field is missing or invalid
func (h *Handler) validateRegisterPayload(payload types.RegisterUserPayload) error {
	// Email validation
	if payload.Email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(payload.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Password validation
	if payload.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(payload.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(payload.Password) > 32 {
		return fmt.Errorf("password must not exceed 32 characters")
	}
	// Check for at least one number and one letter
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(payload.Password)
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(payload.Password)
	if !hasNumber || !hasLetter {
		return fmt.Errorf("password must contain at least one number and one letter")
	}

	// First name validation
	if payload.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if len(payload.FirstName) < 2 {
		return fmt.Errorf("first name must be at least 2 characters long")
	}
	if len(payload.FirstName) > 50 {
		return fmt.Errorf("first name must not exceed 50 characters")
	}
	// Check for valid characters in first name
	if !regexp.MustCompile(`^[a-zA-Z\s-']+$`).MatchString(payload.FirstName) {
		return fmt.Errorf("first name contains invalid characters")
	}

	// Last name validation
	if payload.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if len(payload.LastName) < 2 {
		return fmt.Errorf("last name must be at least 2 characters long")
	}
	if len(payload.LastName) > 50 {
		return fmt.Errorf("last name must not exceed 50 characters")
	}
	// Check for valid characters in last name
	if !regexp.MustCompile(`^[a-zA-Z\s-']+$`).MatchString(payload.LastName) {
		return fmt.Errorf("last name contains invalid characters")
	}

	return nil
}
