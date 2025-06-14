// Package types contains all the shared types and interfaces used across the application
package types

import "time"

// UserStore defines the interface for user data operations
// Any struct that implements these methods can be used as a user store
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
}

type ProductStore interface {
	GetProductByID() ([]Product, error)
	CreateProduct(product *Product) error
}

type Product struct {
	ID          int       `json:"id"`          // Unique identifier for the product
	Name        string    `json:"name"`        // Product name
	Description string    `json:"description"` // Product description
	Image       string    `json:"image"`       // Product image
	Price       float64   `json:"price"`       // Product price
	Quantity    int       `json:"quantity"`    // Product quantity
	CreatedAt   time.Time `json:"createdAt"`   // Timestamp when the product was created
}

type mockUserStore struct {
}

func GetUserByEmail(email string) (*User, error) {
	return nil, nil
}

// User represents a user in the system
// Contains all the user-related fields
type User struct {
	ID        int       `json:"id"`        // Unique identifier for the user
	FirstName string    `json:"firstName"` // User's first name
	LastName  string    `json:"lastName"`  // User's last name
	Email     string    `json:"email"`     // User's email address (unique)
	Password  string    `json:"password"`  // Hashed password
	CreatedAt time.Time `json:"createdAt"` // Timestamp when the user was created
}

// RegisterUserPayload represents the data required for user registration
// Used to validate and process registration requests
type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=30"`         // User's first name
	LastName  string `json:"lastName" validate:"required,min=2,max=30"`          // User's last name
	Email     string `json:"email" validate:"required,email"`                    // User's email address
	Password  string `json:"password" validate:"required,min=8,max=16,alphanum"` // User's password (will be hashed)
}

// LoginUserPayload represents the data required for user login
// Used to validate and process login requests
type LoginUserPayload struct {
	Email    string `json:"email"`    // User's email address
	Password string `json:"password"` // User's password
}
