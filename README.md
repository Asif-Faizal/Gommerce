# Gommerce - Go E-commerce API

A modern e-commerce API built with Go, following clean architecture principles and best practices.

## Project Structure

```xml
.
├── cmd/
│   ├── api/         # API server implementation
│   └── main.go      # Application entry point
├── config/
│   └── env.go       # Environment configuration
├── db/
│   └── db.go        # Database connection setup
├── services/
│   ├── auth/        # Authentication services
│   └── user/        # User-related functionality
├── types/           # Shared types and interfaces
└── utils/           # Utility functions
```

## Current Features

### API Server

- RESTful API server built with Go's standard `net/http` package
- Uses Gorilla Mux for routing
- API versioning support (v1)
- Structured error handling
- MySQL database integration

### Authentication

- Secure password hashing using bcrypt
- User registration with email validation
- Password security best practices
- JWT token support (coming soon)

### User Management

- User registration endpoint
- Email uniqueness validation
- Secure password storage
- User data persistence

### Database

- MySQL database integration
- Connection pooling
- Prepared statements for security
- Error handling and logging

### API Endpoints

#### User Endpoints

##### Register

- **URL**: `/api/v1/register`
- **Method**: `POST`
- **Body**:

  ```json
  {
    "first_name": "string",
    "last_name": "string",
    "email": "string",
    "password": "string"
  }
  ```

- **Status**: Implemented

##### Login

- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Status**: In Progress

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Git
- MySQL Server

### Environment Variables

The following environment variables can be configured (defaults shown):

```xml
PUBLIC_HOST=http://localhost
PORT=8080
DB_USER=root
DB_PASSWORD=root
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=gommerce
```

### Installation

1. Clone the repository:

```bash
git clone https://github.com/Asif-Faizal/Gommerce.git
cd Gommerce
```

2. Install dependencies:

```bash
go mod download
```

3. Set up your environment variables (optional, defaults will be used if not set)

4. Create the database:

```sql
CREATE DATABASE gommerce;
```

5. Run the server:

```bash
go run cmd/main.go
```

The server will start on the configured port (default: `http://localhost:8080`)

## Project Architecture

The project follows a clean, modular architecture:

1. **Entry Point** (`cmd/main.go`)
   - Initializes the database connection
   - Sets up the API server
   - Handles server startup and error logging

2. **API Server** (`cmd/api/api.go`)
   - Manages HTTP server configuration
   - Sets up routing and middleware
   - Handles API versioning

3. **User Service** (`services/user/`)
   - Implements user-related business logic
   - Manages user authentication routes
   - Handles user registration and login
   - Database operations for users

4. **Authentication** (`services/auth/`)
   - Password hashing and validation
   - Security utilities
   - JWT token handling (coming soon)

5. **Types** (`types/`)
   - Shared data structures
   - Interface definitions
   - Request/Response payloads

6. **Utils** (`utils/`)
   - Common utility functions
   - JSON handling
   - Error formatting

## Development Status

The project is currently in active development. Current focus:

- User authentication implementation
- Database schema design
- API endpoint implementation
- Security enhancements

Planned features:

- Product management
- Order processing
- Payment integration
- Admin dashboard

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
