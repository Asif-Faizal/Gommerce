// Package user contains the user-related HTTP handlers and routes
package user

import (
	"net/http"

	"github.com/Asif-Faizal/Gommerce/services/auth"
	"github.com/Asif-Faizal/Gommerce/types"
	"github.com/Asif-Faizal/Gommerce/utils"
	"github.com/gorilla/mux"
)

// Handler represents the user-related HTTP handlers
// It contains methods to handle different user-related endpoints
type Handler struct {
	store types.UserStore
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
	// get JSON payload
	var payload types.RegisterUserPayload
	// validate the request body
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// check if user already exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if does not exist, create user, else return error
	user := &types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	}
	utils.WriteJSON(w, http.StatusOK, user)
	// return the response
}
