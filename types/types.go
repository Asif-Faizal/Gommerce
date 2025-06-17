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
	GetProducts() ([]Product, error)
	CreateProduct(product *Product) error
	GetProductsByIDs(ids []int) ([]Product, error)
}

type OrderStore interface {
	CreateOrder(order *Order) (int, error)
	CreateOrderItem(orderItem *OrderItem) error
	GetOrders(userID int) ([]Order, error)
}

type Order struct {
	ID        int         `json:"id"`        // Unique identifier for the order
	UserID    int         `json:"userID"`    // User ID associated with the order
	Total     float64     `json:"total"`     // Total amount of the order
	Status    string      `json:"status"`    // Status of the order
	Address   string      `json:"address"`   // Address of the order
	CreatedAt time.Time   `json:"createdAt"` // Timestamp when the order was created
	Items     []OrderItem `json:"items"`     // List of items in the order
}

type OrderItem struct {
	ID        int       `json:"id"`        // Unique identifier for the order item
	OrderID   int       `json:"orderID"`   // Order ID associated with the order item
	ProductID int       `json:"productID"` // Product ID associated with the order item
	Quantity  int       `json:"quantity"`  // Quantity of the product in the order item
	Price     float64   `json:"price"`     // Price of the product in the order item
	CreatedAt time.Time `json:"createdAt"` // Timestamp when the order item was created
	Product   *Product  `json:"product"`   // Product details
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

type CartItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}

type CartCheckoutPayload struct {
	Items   []CartItem `json:"items" validate:"required,min=1"`
	Address string     `json:"address" validate:"required"`
}
