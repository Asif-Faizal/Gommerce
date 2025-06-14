// Package user contains tests for user-related HTTP handlers and routes
package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Asif-Faizal/Gommerce/services/auth"
	"github.com/Asif-Faizal/Gommerce/types"
	"github.com/gorilla/mux"
)

// TestUserServiceHandlers is the main test function for user service handlers
// It sets up the test environment and runs individual test cases
func TestUserServiceHandlers(t *testing.T) {
	// Create a mock user store for testing
	userStore := &mockUserStore{}
	// Create a new handler with the mock store
	handler := NewHandler(userStore)

	// Test case: Invalid user registration payload
	t.Run("Should fail if payload is invalid", func(t *testing.T) {
		testCases := []struct {
			name    string
			payload types.RegisterUserPayload
			wantErr string
		}{
			// Email validation cases
			{
				name: "empty email",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "",
					Password:  "password123",
				},
				wantErr: "email is required",
			},
			{
				name: "invalid email format",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "invalid-email",
					Password:  "password123",
				},
				wantErr: "invalid email format",
			},
			{
				name: "missing @ in email",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "testexample.com",
					Password:  "password123",
				},
				wantErr: "invalid email format",
			},
			{
				name: "missing domain in email",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@",
					Password:  "password123",
				},
				wantErr: "invalid email format",
			},

			// Password validation cases
			{
				name: "empty password",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "",
				},
				wantErr: "password is required",
			},
			{
				name: "password too short",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "pass",
				},
				wantErr: "password must be at least 8 characters long",
			},
			{
				name: "password too long",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "thispasswordiswaytoolongandshouldnotbeaccepted",
				},
				wantErr: "password must not exceed 32 characters",
			},
			{
				name: "password without numbers",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "passwordonly",
				},
				wantErr: "password must contain at least one number and one letter",
			},
			{
				name: "password without letters",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "12345678",
				},
				wantErr: "password must contain at least one number and one letter",
			},

			// First name validation cases
			{
				name: "empty first name",
				payload: types.RegisterUserPayload{
					FirstName: "",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "first name is required",
			},
			{
				name: "first name too short",
				payload: types.RegisterUserPayload{
					FirstName: "J",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "first name must be at least 2 characters long",
			},
			{
				name: "first name too long",
				payload: types.RegisterUserPayload{
					FirstName: "ThisFirstNameIsWayTooLongAndShouldNotBeAcceptedInTheSystem",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "first name must not exceed 50 characters",
			},
			{
				name: "first name with invalid characters",
				payload: types.RegisterUserPayload{
					FirstName: "John123",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "first name contains invalid characters",
			},

			// Last name validation cases
			{
				name: "empty last name",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "last name is required",
			},
			{
				name: "last name too short",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "D",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "last name must be at least 2 characters long",
			},
			{
				name: "last name too long",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "ThisLastNameIsWayTooLongAndShouldNotBeAcceptedInTheSystem",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "last name must not exceed 50 characters",
			},
			{
				name: "last name with invalid characters",
				payload: types.RegisterUserPayload{
					FirstName: "John",
					LastName:  "Doe123",
					Email:     "test@example.com",
					Password:  "password123",
				},
				wantErr: "last name contains invalid characters",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Marshal the payload to JSON
				marshaled, err := json.Marshal(tc.payload)
				if err != nil {
					t.Fatalf("Failed to marshal payload: %v", err)
				}

				// Create a new HTTP request with the JSON payload
				req, err := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(marshaled))
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}

				// Create a response recorder to capture the response
				rr := httptest.NewRecorder()

				// Set up the router and register the handler
				router := mux.NewRouter()
				router.HandleFunc("/user/register", handler.handleRegister).Methods(http.MethodPost)

				// Serve the request
				router.ServeHTTP(rr, req)

				// Check if the response status code is 400 (Bad Request)
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
	t.Run("Should fail if user already exists", func(t *testing.T) {
		// Create a mock store that returns an existing user
		mockStore := &mockUserStore{
			getUserByEmailFunc: func(email string) (*types.User, error) {
				return &types.User{
					ID:        1,
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Password:  "hashedpassword",
					CreatedAt: time.Now(),
				}, nil
			},
		}

		handler := NewHandler(mockStore)
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "test@example.com",
			Password:  "password123",
		}

		marshaled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user/register", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		var response map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedErr := fmt.Sprintf("user with email %s already exists", payload.Email)
		if response["error"] != expectedErr {
			t.Errorf("Expected error %q, got %q", expectedErr, response["error"])
		}
	})
	t.Run("Should create a new user if payload is valid", func(t *testing.T) {
		// Create a valid payload
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "test@example.com",
			Password:  "password123",
		}

		// Set up mock functions
		getUserByEmailCalled := false
		createUserCalled := false

		userStore.getUserByEmailFunc = func(email string) (*types.User, error) {
			getUserByEmailCalled = true
			return nil, sql.ErrNoRows
		}

		userStore.createUserFunc = func(user *types.User) error {
			createUserCalled = true
			user.ID = 1 // Set an ID for the created user
			return nil
		}

		// Marshal the payload to JSON
		marshaled, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}

		// Create a new HTTP request with the JSON payload
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Set up the router and register the handler
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)

		// Serve the request
		router.ServeHTTP(rr, req)

		// Check if the response status code is 201 (Created)
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Verify success message
		if response["message"] != "user created successfully" {
			t.Errorf("Expected message %q, got %q", "user created successfully", response["message"])
		}

		// Verify function calls
		if !getUserByEmailCalled {
			t.Error("GetUserByEmail was not called")
		}
		if !createUserCalled {
			t.Error("CreateUser was not called")
		}
	})

	// Test login functionality
	t.Run("Login Tests", func(t *testing.T) {
		// Create a test password and hash it
		testPassword := "password123"
		hashedPassword, err := auth.HashPassword(testPassword)
		if err != nil {
			t.Fatalf("Failed to hash test password: %v", err)
		}

		testCases := []struct {
			name          string
			payload       types.LoginUserPayload
			mockUser      *types.User
			mockError     error
			expectedCode  int
			expectedError string
		}{
			{
				name: "successful login",
				payload: types.LoginUserPayload{
					Email:    "test@example.com",
					Password: testPassword,
				},
				mockUser: &types.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  hashedPassword,
				},
				expectedCode: http.StatusOK,
			},
			{
				name: "user not found",
				payload: types.LoginUserPayload{
					Email:    "nonexistent@example.com",
					Password: testPassword,
				},
				mockError:     sql.ErrNoRows,
				expectedCode:  http.StatusUnauthorized,
				expectedError: "invalid email or password",
			},
			{
				name: "invalid password",
				payload: types.LoginUserPayload{
					Email:    "test@example.com",
					Password: "wrongpassword",
				},
				mockUser: &types.User{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "test@example.com",
					Password:  hashedPassword,
				},
				expectedCode:  http.StatusUnauthorized,
				expectedError: "invalid email or password",
			},
			{
				name: "empty email",
				payload: types.LoginUserPayload{
					Email:    "",
					Password: testPassword,
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "email is required",
			},
			{
				name: "invalid email format",
				payload: types.LoginUserPayload{
					Email:    "invalid-email",
					Password: testPassword,
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "invalid email format",
			},
			{
				name: "empty password",
				payload: types.LoginUserPayload{
					Email:    "test@example.com",
					Password: "",
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "password is required",
			},
			{
				name: "password too short",
				payload: types.LoginUserPayload{
					Email:    "test@example.com",
					Password: "short",
				},
				expectedCode:  http.StatusBadRequest,
				expectedError: "password must be at least 8 characters long",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create mock store
				mockStore := &mockUserStore{
					getUserByEmailFunc: func(email string) (*types.User, error) {
						if tc.mockError != nil {
							return nil, tc.mockError
						}
						return tc.mockUser, nil
					},
				}

				handler := NewHandler(mockStore)

				// Create request
				payload, err := json.Marshal(tc.payload)
				if err != nil {
					t.Fatalf("Failed to marshal payload: %v", err)
				}

				req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}

				// Create response recorder
				rr := httptest.NewRecorder()

				// Create router and register handler
				router := mux.NewRouter()
				router.HandleFunc("/login", handler.handleLogin).Methods(http.MethodPost)

				// Serve request
				router.ServeHTTP(rr, req)

				// Check status code
				if rr.Code != tc.expectedCode {
					t.Errorf("Expected status %d, got %d", tc.expectedCode, rr.Code)
				}

				// Check response body
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
					if response["message"] != "login successful" {
						t.Error("Expected message 'login successful' in response")
					}
					if data, ok := response["data"].(map[string]interface{}); ok {
						if user, ok := data["user"].(map[string]interface{}); ok {
							if user["email"] != tc.payload.Email {
								t.Error("Email mismatch in response data")
							}
						} else {
							t.Error("Expected user object in response data")
						}
					} else {
						t.Error("Expected data object in response")
					}
				}
			})
		}
	})
}

// mockUserStore implements the types.UserStore interface for testing
type mockUserStore struct {
	getUserByEmailFunc func(email string) (*types.User, error)
	createUserFunc     func(user *types.User) error
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(email)
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user *types.User) error {
	if m.createUserFunc != nil {
		return m.createUserFunc(user)
	}
	return nil
}
