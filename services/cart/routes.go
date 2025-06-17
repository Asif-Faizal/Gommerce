package cart

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Asif-Faizal/Gommerce/types"
	"github.com/Asif-Faizal/Gommerce/utils"
	"github.com/gorilla/mux"
)

// Handler represents the user-related HTTP handlers
// It contains methods to handle different user-related endpoints
type Handler struct {
	store        types.OrderStore   // Interface for user data operations
	productStore types.ProductStore // Interface for product data operations
}

// NewHandler creates a new instance of the user Handler
// This is a constructor function for the Handler struct
func NewHandler(store types.OrderStore, productStore types.ProductStore) *Handler {
	return &Handler{store: store, productStore: productStore}
}

func (h *Handler) OrderRoutes(router *mux.Router) {
	router.HandleFunc("/order", h.handleCheckout).Methods(http.MethodPost)
	router.HandleFunc("/orders", h.handleGetOrders).Methods(http.MethodGet)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.AuthenticateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.Validate.Struct(cart); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get products
	productIDs := make([]int, len(cart.Items))
	for i, item := range cart.Items {
		productIDs[i] = item.ProductID
	}
	products, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validate products exist
	if len(products) != len(productIDs) {
		http.Error(w, "one or more products not found", http.StatusBadRequest)
		return
	}

	// calculate total and validate quantities
	total := 0.0
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	for _, item := range cart.Items {
		product, exists := productMap[item.ProductID]
		if !exists {
			http.Error(w, fmt.Sprintf("product with ID %d not found", item.ProductID), http.StatusBadRequest)
			return
		}
		if item.Quantity > product.Quantity {
			http.Error(w, fmt.Sprintf("insufficient quantity for product %d", item.ProductID), http.StatusBadRequest)
			return
		}
		total += product.Price * float64(item.Quantity)
	}

	// create order
	order := &types.Order{
		UserID:    userId,
		Total:     total,
		Status:    "pending",
		Address:   cart.Address,
		CreatedAt: time.Now(),
	}

	// create order in database
	orderID, err := h.store.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	order.ID = orderID

	// create order items
	for _, item := range cart.Items {
		product := productMap[item.ProductID]
		orderItem := &types.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		if err := h.store.CreateOrderItem(orderItem); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// return success response
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "order created successfully",
		"data":    order,
	})
}

func (h *Handler) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.AuthenticateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	orders, err := h.store.GetOrders(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "orders fetched successfully",
		"data":    orders,
	})
}
