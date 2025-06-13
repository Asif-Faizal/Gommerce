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
└── services/
    └── user/        # User-related functionality
```

## Current Features

### API Server

- RESTful API server built with Go's standard `net/http` package
- Uses Gorilla Mux for routing
- API versioning support (v1)
- Structured error handling
- MySQL database integration

### Configuration Management

- Environment-based configuration
- Default values for development
- Configurable database settings
- Flexible port and host configuration

### Database

- MySQL database integration
- Configurable connection settings
- Connection pooling support
- Environment-based configuration

### User Service

- User authentication endpoints:
  - `/api/v1/login` - User login
  - `/api/v1/register` - User registration
- Modular handler structure for easy maintenance
- Clean separation of concerns

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
DB_HOST=http://127.0.0.1
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

4. Run the server:

```bash
go run cmd/main.go
```

The server will start on the configured port (default: `http://localhost:8080`)

## API Endpoints

### User Endpoints

#### Login

- **URL**: `/api/v1/login`
- **Method**: `POST`
- **Status**: To be implemented

#### Register

- **URL**: `/api/v1/register`
- **Method**: `POST`
- **Status**: To be implemented

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

3. **Configuration** (`config/env.go`)
   - Manages environment variables
   - Provides default configurations
   - Handles database settings

4. **Database** (`db/db.go`)
   - Manages database connections
   - Handles connection pooling
   - Provides database utilities

5. **User Service** (`services/user/routes.go`)
   - Implements user-related business logic
   - Manages user authentication routes
   - Handles user registration and login

## Development Status

The project is currently in early development. The following features are planned:

- User authentication implementation
- Product management
- Order processing
- Payment integration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
