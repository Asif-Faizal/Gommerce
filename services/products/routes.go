package products

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Asif-Faizal/Gommerce/config"
	"github.com/Asif-Faizal/Gommerce/services/auth"
	"github.com/Asif-Faizal/Gommerce/types"
	"github.com/Asif-Faizal/Gommerce/utils"
	"github.com/gorilla/mux"
)

// Handler represents the user-related HTTP handlers
// It contains methods to handle different user-related endpoints
type Handler struct {
	store types.ProductStore // Interface for user data operations
}

// NewHandler creates a new instance of the user Handler
// This is a constructor function for the Handler struct
func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

// authenticateRequest is a helper function to authenticate requests
func authenticateRequest(r *http.Request) (int, error) {
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

// RegisterRoutes sets up all the user-related routes
// It takes a router and attaches the handler functions to specific paths
func (h *Handler) ProductRoutes(router *mux.Router) {
	router.HandleFunc("/products/create", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request
	userId, err := authenticateRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	log.Printf("User %d requesting products list", userId)

	products, err := h.store.GetProductByID()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "products fetched successfully",
		"data":    products,
	})
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// Authenticate the request
	userId, err := authenticateRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	log.Printf("User %d attempting to create a product", userId)

	var product types.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Decoded product: %+v", product)

	// Validate required fields
	if product.Name == "" {
		log.Printf("Validation error: name is required")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("name is required"))
		return
	}
	if product.Description == "" {
		log.Printf("Validation error: description is required")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("description is required"))
		return
	}
	if product.Image == "" {
		log.Printf("Validation error: image is required")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image is required"))
		return
	}
	if product.Price <= 0 {
		log.Printf("Validation error: price must be greater than 0")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("price must be greater than 0"))
		return
	}
	if product.Quantity < 0 {
		log.Printf("Validation error: quantity cannot be negative")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("quantity cannot be negative"))
		return
	}

	log.Printf("Creating product in database")
	if err := h.store.CreateProduct(&product); err != nil {
		log.Printf("Error creating product: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Product created successfully with ID: %d by user: %d", product.ID, userId)
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "product created successfully",
		"data":    product,
	})
}
