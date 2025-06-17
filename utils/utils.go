// Package utils contains common utility functions used across the application
package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Asif-Faizal/Gommerce/config"
	"github.com/Asif-Faizal/Gommerce/services/auth"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// ParseJSON parses the JSON body of an HTTP request into the provided payload
// Returns an error if the body is nil or if JSON parsing fails
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is nil")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJSON writes a JSON response to the HTTP response writer
// Sets the content type to application/json and the provided status code
// Returns any potential error during JSON encoding
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// WriteError writes an error response to the HTTP response writer
// Formats the error message into a JSON response with the provided status code
// Returns any potential error during JSON encoding
func WriteError(w http.ResponseWriter, status int, err error) error {
	return WriteJSON(w, status, map[string]string{"error": err.Error()})
}

// authenticateRequest is a helper function to authenticate requests
func AuthenticateRequest(r *http.Request) (int, error) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("authorization header is required")
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, fmt.Errorf("invalid authorization header format")
	}

	// Verify the token
	secret := []byte(config.Envs.JWTSecret)
	userId, err := auth.VerifyJWT(parts[1], secret)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	return userId, nil
}
