package products

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

// RegisterRoutes sets up all the user-related routes
// It takes a router and attaches the handler functions to specific paths
func (h *Handler) ProductRoutes(router *mux.Router) {
	router.HandleFunc("/products/create", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("Received create product request")

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

	log.Printf("Product created successfully with ID: %d", product.ID)
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "product created successfully",
		"data":    product,
	})
}
