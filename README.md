# Gommerce - Go E-commerce API

A modern e-commerce API built with Go, following clean architecture principles and best practices. This project is designed to be both educational and production-ready, making it perfect for developers learning Go or building their first e-commerce application.

## Why This Project?

This project was created to demonstrate:
- How to build a production-grade Go API
- Best practices in Go development
- Clean architecture principles
- Test-driven development (TDD)
- Secure authentication practices
- Database integration patterns

## Project Structure Explained

```xml
.
├── cmd/                    # Command-line applications
│   ├── api/               # API server implementation
│   │   ├── api.go        # Server setup and configuration
│   │   └── routes.go     # Route definitions
│   └── main.go           # Application entry point
├── config/                # Configuration management
│   └── env.go            # Environment variable handling
├── db/                    # Database layer
│   └── db.go             # Database connection and queries
├── services/             # Business logic layer
│   ├── auth/             # Authentication services
│   │   └── auth.go       # Auth-related business logic
│   └── user/             # User management
│       ├── routes.go     # User-related HTTP handlers
│       └── routes_test.go # Tests for user handlers
├── types/                # Shared types and interfaces
│   └── types.go          # Common data structures
└── utils/                # Utility functions
    └── utils.go          # Helper functions
```

### Key Components Explained

1. **cmd/** - Contains the main application code
   - `main.go`: The starting point of our application
   - `api/`: Contains the HTTP server setup and routing

2. **services/** - Contains our business logic
   - Each service (user, auth) is isolated in its own package
   - Includes both implementation and tests
   - Follows the single responsibility principle

3. **types/** - Defines our data structures
   - Shared interfaces and structs
   - Clear separation of concerns
   - Makes the code more maintainable

4. **utils/** - Common helper functions
   - Reusable across the application
   - Keeps our code DRY (Don't Repeat Yourself)

## Development Approach

### 1. Testing Strategy
We follow Test-Driven Development (TDD):
- Write tests first
- Implement the feature
- Refactor while keeping tests passing
- Example from our codebase:
  ```go
  // routes_test.go
  func TestUserServiceHandlers(t *testing.T) {
      // Tests user registration validation
      // Ensures empty emails are rejected
      // Verifies proper error responses
  }
  ```

### 2. Error Handling
We implement comprehensive error handling:
- Clear error messages
- Proper HTTP status codes
- Consistent error response format
- Example:
  ```go
  // Returns 400 for invalid input
  // Returns 500 for server errors
  // Returns 404 for not found
  ```

### 3. Security Practices
Security is a top priority:
- Password hashing using bcrypt
- Input validation
- SQL injection prevention
- Secure session management

## Current Features

### API Server
- RESTful API using Go's `net/http`
- Gorilla Mux for routing
- API versioning (v1)
- Structured error handling

### Authentication
- Secure password hashing
- User registration with validation
- JWT token support (coming soon)

### User Management
- Registration endpoint
- Email validation
- Secure password storage

### Database
- MySQL integration
- Connection pooling
- Prepared statements
- Error handling

## API Endpoints

### User Registration
```http
POST /api/v1/register
Content-Type: application/json

{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "password": "securepassword"
}
```

Response:
```json
{
    "status": "success",
    "message": "User registered successfully"
}
```

### User Login
```http
POST /api/v1/login
Content-Type: application/json

{
    "email": "john@example.com",
    "password": "securepassword"
}
```

Response:
```json
{
    "status": "success",
    "message": "Login successful",
    "data": {
        "id": 1,
        "firstName": "John",
        "lastName": "Doe",
        "email": "john@example.com"
    }
}
```

### Get User by ID
```http
GET /api/v1/users/{id}
Authorization: Bearer {token}
```

Response:
```json
{
    "status": "success",
    "data": {
        "id": 1,
        "firstName": "John",
        "lastName": "Doe",
        "email": "john@example.com"
    }
}
```

### Curl Commands

1. **Register a new user**:
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

2. **Login with user credentials**:
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

3. **Get user details** (requires authentication):
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Response Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

### Error Response Format
```json
{
    "error": "Error message description"
}
```

## Implementation Deep Dive: User Registration

Let's walk through the complete implementation of the user registration feature, from design to testing to implementation.

### 1. Design Phase

#### Requirements Analysis
- User needs to provide: first name, last name, email, and password
- Email must be unique
- Password must be secure
- All fields are required
- Response should be consistent

#### Data Structures
```go
// types/types.go
type RegisterUserPayload struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}

type User struct {
    ID        int       `json:"id"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    CreatedAt time.Time `json:"createdAt"`
}
```

### 2. Testing Phase (TDD)

#### Test Cases
```go
// services/user/routes_test.go
func TestUserServiceHandlers(t *testing.T) {
    // Test Case 1: Invalid Payload
    t.Run("Should fail if user payload is invalid", func(t *testing.T) {
        payload := types.RegisterUserPayload{
            FirstName: "John",
            LastName:  "Doe",
            Email:     "", // Empty email should fail
            Password:  "password",
        }
        // ... test implementation
    })

    // Test Case 2: Duplicate Email
    t.Run("Should fail if email already exists", func(t *testing.T) {
        // ... test implementation
    })

    // Test Case 3: Successful Registration
    t.Run("Should successfully register a new user", func(t *testing.T) {
        // ... test implementation
    })
}
```

#### Mock Implementation
```go
// services/user/routes_test.go
type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
    return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) CreateUser(user *types.User) error {
    return nil
}
```

### 3. Implementation Phase

#### Handler Function
```go
// services/user/routes.go
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
    // 1. Parse JSON payload
    var payload types.RegisterUserPayload
    if err := utils.ParseJSON(r, &payload); err != nil {
        utils.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    // 2. Validate payload
    if err := validateRegisterPayload(payload); err != nil {
        utils.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    // 3. Check if user exists
    user, err := h.store.GetUserByEmail(payload.Email)
    if err == nil && user != nil {
        utils.WriteError(w, http.StatusBadRequest, "user with that email already exists")
        return
    }

    // 4. Hash password
    hashedPassword, err := auth.HashPassword(payload.Password)
    if err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "error creating user")
        return
    }

    // 5. Create user
    user = &types.User{
        FirstName: payload.FirstName,
        LastName:  payload.LastName,
        Email:     payload.Email,
        Password:  hashedPassword,
        CreatedAt: time.Now(),
    }

    if err := h.store.CreateUser(user); err != nil {
        utils.WriteError(w, http.StatusInternalServerError, "error creating user")
        return
    }

    // 6. Return success response
    utils.WriteJSON(w, http.StatusCreated, map[string]string{
        "message": "user created successfully",
    })
}
```

#### Validation Function
```go
// services/user/routes.go
func validateRegisterPayload(payload types.RegisterUserPayload) error {
    if payload.Email == "" {
        return fmt.Errorf("email is required")
    }
    if payload.Password == "" {
        return fmt.Errorf("password is required")
    }
    if payload.FirstName == "" {
        return fmt.Errorf("first name is required")
    }
    if payload.LastName == "" {
        return fmt.Errorf("last name is required")
    }
    return nil
}
```

### 4. Security Considerations

1. **Password Security**
   - Passwords are hashed using bcrypt
   - Never store plain text passwords
   - Use secure hashing algorithm

2. **Input Validation**
   - Validate all input fields
   - Sanitize user input
   - Prevent SQL injection

3. **Error Handling**
   - Don't expose internal errors
   - Use appropriate HTTP status codes
   - Provide clear error messages

### 5. API Usage Example

#### Request
```http
POST /api/v1/register
Content-Type: application/json

{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "password": "securepassword123"
}
```

#### Success Response
```json
{
    "message": "user created successfully"
}
```

#### Error Response
```json
{
    "error": "email is required"
}
```

### 6. Testing the Implementation

1. **Run the Tests**
```bash
go test ./services/user -v
```

2. **Manual Testing**
```bash
# Using curl
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### 7. Best Practices Implemented

1. **Code Organization**
   - Separation of concerns
   - Clear function responsibilities
   - Modular design

2. **Error Handling**
   - Consistent error responses
   - Proper HTTP status codes
   - Clear error messages

3. **Testing**
   - Unit tests for all cases
   - Mock implementations
   - Edge case coverage

4. **Security**
   - Password hashing
   - Input validation
   - Error message sanitization

This implementation follows Go best practices and provides a solid foundation for user registration in our e-commerce application.

## Getting Started

### Prerequisites
- Go 1.16+
- Git
- MySQL Server

### Environment Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/Asif-Faizal/Gommerce.git
   cd Gommerce
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Configure environment variables:
   ```bash
   # Copy the example env file
   cp .env.example .env
   
   # Edit .env with your settings
   ```

4. Create the database:
   ```sql
   CREATE DATABASE gommerce;
   ```

5. Run the server:
   ```bash
   go run cmd/main.go
   ```

### Running Tests
```bash
# Run all tests
go test ./...

# Run specific test
go test ./services/user -v
```

## Development Workflow

1. **Write Tests First**
   - Define expected behavior
   - Write test cases
   - Run tests (they should fail)

2. **Implement Features**
   - Write the minimum code to pass tests
   - Follow Go best practices
   - Keep code clean and documented

3. **Review and Refactor**
   - Check code quality
   - Improve performance
   - Update documentation

## Contributing

We welcome contributions! Here's how:

1. Fork the repository
2. Create a feature branch
3. Write tests for your feature
4. Implement the feature
5. Submit a pull request

## Learning Resources

- [Go Documentation](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Go Testing](https://golang.org/pkg/testing/)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Database Migrations

### Migration Structure
```xml
cmd/migrate/
├── migrations/           # Migration files
│   ├── 000001_add-user-table.up.sql   # Creates users table
│   └── 000001_add-user-table.down.sql # Drops users table
└── main.go              # Migration runner
```

### Users Table Schema
```sql
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT,
    firstName VARCHAR(50) NOT NULL,
    lastName VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Migration Runner Implementation

The migration runner (`cmd/migrate/main.go`) handles database migrations using the `golang-migrate` library. Here's how it works:

```go
// Main components of the migration runner
1. Database Connection
   - Uses MySQL driver
   - Configures connection using environment variables
   - Enables native passwords and time parsing

2. Migration Driver
   - Creates MySQL-specific migration driver
   - Handles database-specific operations
   - Manages transaction safety

3. Migration Instance
   - Points to migrations directory
   - Supports both up and down migrations
   - Handles migration versioning
```

#### Environment Configuration
The migration runner uses the following environment variables:
```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_ADDRESS=localhost:3306
DB_NAME=gommerce
```

#### Command Line Usage
The migration runner accepts two commands:
```bash
# Apply migrations
go run cmd/migrate/main.go up

# Rollback migrations
go run cmd/migrate/main.go down
```

#### Error Handling
The migration runner includes robust error handling:
- Database connection errors
- Migration driver initialization errors
- Migration execution errors
- Version control errors

#### Logging
The runner provides detailed logging:
- Configuration details
- Database connection status
- Migration execution results
- Error messages when applicable

### Migration Workflow

1. **Initial Setup**
   ```bash
   # Create database
   mysql -u root -p
   CREATE DATABASE gommerce;
   
   # Run initial migration
   make migrate-up
   ```

2. **Creating New Migrations**
   ```bash
   # Create a new migration
   make migration-create add_new_table
   
   # This creates:
   # - cmd/migrate/migrations/XXXXXX_add_new_table.up.sql
   # - cmd/migrate/migrations/XXXXXX_add_new_table.down.sql
   ```

3. **Applying Migrations**
   ```bash
   # Apply all pending migrations
   make migrate-up
   
   # Rollback last migration
   make migrate-down
   ```

4. **Verifying Migrations**
   ```bash
   # Check migration status
   mysql -u root -p gommerce
   SELECT * FROM schema_migrations;
   ```

### Migration Best Practices (Updated)

1. **Version Control**
   - Always commit both up and down migrations
   - Never modify existing migrations
   - Create new migrations for changes
   - Keep migrations in version control

2. **Naming Conventions**
   - Use descriptive names
   - Follow pattern: `NNNNNN_description.up.sql`
   - Include corresponding down migration
   - Use snake_case for table and column names

3. **Data Integrity**
   - Include proper constraints
   - Use appropriate data types
   - Add necessary indexes
   - Handle foreign keys properly

4. **Error Handling**
   - Test migrations before deployment
   - Include rollback procedures
   - Handle edge cases
   - Log migration results

5. **Performance**
   - Use appropriate indexes
   - Optimize large migrations
   - Consider batch operations
   - Monitor migration execution time
