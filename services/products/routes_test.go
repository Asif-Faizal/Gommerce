package products

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Asif-Faizal/Gommerce/types"
	"github.com/gorilla/mux"
)

// TestProductServiceHandlers is the main test function for product service handlers
func TestProductServiceHandlers(t *testing.T) {
	// Create a mock product store for testing
	productStore := &mockProductStore{}
	// Create a new handler with the mock store
	handler := NewHandler(productStore)

	// Test case: Create Product Validation
	t.Run("Should fail if create product payload is invalid", func(t *testing.T) {
		testCases := []struct {
			name    string
			payload types.Product
			wantErr string
		}{
			{
				name: "empty name",
				payload: types.Product{
					Description: "Test Description",
					Image:       "https://example.com/image.jpg",
					Price:       99.99,
					Quantity:    10,
				},
				wantErr: "name is required",
			},
			{
				name: "empty description",
				payload: types.Product{
					Name:     "Test Product",
					Image:    "https://example.com/image.jpg",
					Price:    99.99,
					Quantity: 10,
				},
				wantErr: "description is required",
			},
			{
				name: "empty image",
				payload: types.Product{
					Name:        "Test Product",
					Description: "Test Description",
					Price:       99.99,
					Quantity:    10,
				},
				wantErr: "image is required",
			},
			{
				name: "invalid price",
				payload: types.Product{
					Name:        "Test Product",
					Description: "Test Description",
					Image:       "https://example.com/image.jpg",
					Price:       0,
					Quantity:    10,
				},
				wantErr: "price must be greater than 0",
			},
			{
				name: "negative price",
				payload: types.Product{
					Name:        "Test Product",
					Description: "Test Description",
					Image:       "https://example.com/image.jpg",
					Price:       -10.0,
					Quantity:    10,
				},
				wantErr: "price must be greater than 0",
			},
			{
				name: "negative quantity",
				payload: types.Product{
					Name:        "Test Product",
					Description: "Test Description",
					Image:       "https://example.com/image.jpg",
					Price:       99.99,
					Quantity:    -1,
				},
				wantErr: "quantity cannot be negative",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Marshal the payload to JSON
				marshaled, err := json.Marshal(tc.payload)
				if err != nil {
					t.Fatalf("Failed to marshal payload: %v", err)
				}

				// Create request
				req, err := http.NewRequest(http.MethodPost, "/products/create", bytes.NewBuffer(marshaled))
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}

				// Create response recorder
				rr := httptest.NewRecorder()

				// Create router and register handler
				router := mux.NewRouter()
				router.HandleFunc("/products/create", handler.handleCreateProduct).Methods(http.MethodPost)

				// Serve request
				router.ServeHTTP(rr, req)

				// Check status code
				if rr.Code != http.StatusBadRequest {
					t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
				}

				var response map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response["error"] != tc.wantErr {
					t.Errorf("Expected error %q, got %q", tc.wantErr, response["error"])
				}
			})
		}
	})

	// Test case: Successful Product Creation
	t.Run("Should create a new product if payload is valid", func(t *testing.T) {
		// Create a valid payload
		payload := types.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Image:       "https://example.com/image.jpg",
			Price:       99.99,
			Quantity:    10,
		}

		// Set up mock function
		createProductCalled := false
		productStore.createProductFunc = func(product *types.Product) error {
			createProductCalled = true
			product.ID = 1 // Set an ID for the created product
			product.CreatedAt = time.Now()
			return nil
		}

		// Marshal the payload to JSON
		marshaled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}

		// Create request
		req, err := http.NewRequest(http.MethodPost, "/products/create", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Create response recorder
		rr := httptest.NewRecorder()

		// Create router and register handler
		router := mux.NewRouter()
		router.HandleFunc("/products/create", handler.handleCreateProduct).Methods(http.MethodPost)

		// Serve request
		router.ServeHTTP(rr, req)

		// Check status code
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify success message
		if response["message"] != "product created successfully" {
			t.Errorf("Expected message %q, got %q", "product created successfully", response["message"])
		}

		// Verify function calls
		if !createProductCalled {
			t.Error("CreateProduct was not called")
		}
	})

	// Test case: Get Products
	t.Run("Get Products Tests", func(t *testing.T) {
		testCases := []struct {
			name          string
			mockProducts  []types.Product
			mockError     error
			expectedCode  int
			expectedError string
		}{
			{
				name: "successful products fetch",
				mockProducts: []types.Product{
					{
						ID:          1,
						Name:        "Product 1",
						Description: "Description 1",
						Image:       "https://example.com/image1.jpg",
						Price:       99.99,
						Quantity:    10,
						CreatedAt:   time.Now(),
					},
					{
						ID:          2,
						Name:        "Product 2",
						Description: "Description 2",
						Image:       "https://example.com/image2.jpg",
						Price:       149.99,
						Quantity:    5,
						CreatedAt:   time.Now(),
					},
				},
				expectedCode: http.StatusOK,
			},
			{
				name:          "database error",
				mockError:     sql.ErrConnDone,
				expectedCode:  http.StatusInternalServerError,
				expectedError: "sql: connection is already closed",
			},
			{
				name:         "empty products list",
				mockProducts: []types.Product{},
				expectedCode: http.StatusOK,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create mock store
				mockStore := &mockProductStore{
					getProductsFunc: func() ([]types.Product, error) {
						if tc.mockError != nil {
							return nil, tc.mockError
						}
						return tc.mockProducts, nil
					},
				}

				handler := NewHandler(mockStore)

				// Create request
				req, err := http.NewRequest(http.MethodGet, "/products", nil)
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}

				// Create response recorder
				rr := httptest.NewRecorder()

				// Create router and register handler
				router := mux.NewRouter()
				router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)

				// Serve request
				router.ServeHTTP(rr, req)

				// Check status code
				if rr.Code != tc.expectedCode {
					t.Errorf("Expected status %d, got %d", tc.expectedCode, rr.Code)
				}

				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if tc.expectedError != "" {
					if response["error"] != tc.expectedError {
						t.Errorf("Expected error %q, got %q", tc.expectedError, response["error"])
					}
				} else {
					// Check success response structure
					if response["status"] != "success" {
						t.Error("Expected status 'success' in response")
					}
					if response["message"] != "products fetched successfully" {
						t.Error("Expected message 'products fetched successfully' in response")
					}

					// Check data
					if data, ok := response["data"].([]interface{}); ok {
						if len(data) != len(tc.mockProducts) {
							t.Errorf("Expected %d products, got %d", len(tc.mockProducts), len(data))
						}
					} else {
						t.Error("Expected data to be an array of products")
					}
				}
			})
		}
	})
}

// mockProductStore implements the types.ProductStore interface for testing
type mockProductStore struct {
	getProductsFunc      func() ([]types.Product, error)
	createProductFunc    func(product *types.Product) error
	getProductsByIDsFunc func(ids []int) ([]types.Product, error)
}

func (m *mockProductStore) GetProducts() ([]types.Product, error) {
	if m.getProductsFunc != nil {
		return m.getProductsFunc()
	}
	return nil, fmt.Errorf("products not found")
}

func (m *mockProductStore) CreateProduct(product *types.Product) error {
	if m.createProductFunc != nil {
		return m.createProductFunc(product)
	}
	return nil
}

func (m *mockProductStore) GetProductsByIDs(ids []int) ([]types.Product, error) {
	if m.getProductsByIDsFunc != nil {
		return m.getProductsByIDsFunc(ids)
	}
	return nil, fmt.Errorf("products not found")
}
