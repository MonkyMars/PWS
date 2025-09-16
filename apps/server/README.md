# PWS Server

The backend API server for the PWS (Personal Web Space) application, built with Go and the Fiber web framework. This server provides a high-performance, scalable REST API with structured logging, comprehensive configuration management, and standardized response handling.

## Overview

This Go application provides:

- High-performance HTTP server using Fiber v3
- Structured logging with slog
- Environment-based configuration management
- Standardized response handling with pagination support
- CORS middleware for cross-origin requests
- Type-safe request/response handling
- Extensible middleware architecture

## Technology Stack

- **Go 1.25.0** - Programming language
- **Fiber v3.0.0-rc.1** - Web framework for high performance
- **slog** - Structured logging (Go standard library)
- **goccy/go-json** - High-performance JSON encoding
- **godotenv** - Environment variable loading

## Project Structure

```
apps/server/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── api/                    # API layer
│   ├── router.go          # Main router configuration
│   ├── config/            # Configuration management
│   │   ├── env.go         # Environment variables
│   │   ├── logger.go      # Logging configuration
│   │   ├── fiber.go       # Fiber server setup
│   │   ├── env_test.go    # Configuration tests
│   │   └── README.md      # Config documentation
│   ├── middleware/        # HTTP middleware
│   │   ├── cors.go        # CORS middleware
│   │   └── README.md      # Middleware documentation
│   ├── response/          # Response handling
│   │   ├── types.go       # Response type definitions
│   │   ├── success.go     # Success response builders
│   │   ├── error.go       # Error response builders
│   │   ├── lib.go         # Response utilities
│   │   └── README.md      # Response documentation
│   └── routes/            # Route handlers
│       └── README.md      # Routes documentation
├── services/              # Business logic services
│   └── database.go        # Database service
└── types/                 # Shared type definitions
    └── README.md          # Types documentation
```

## Getting Started

### Prerequisites

- **Go 1.25.0** - [Download Go](https://golang.org/dl/)
- **Git** - For version control

### Installation

1. Clone the repository and navigate to the server directory:

   ```bash
   cd apps/server
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Copy environment configuration:

   ```bash
   cp .env.example .env
   ```

4. Configure environment variables in `.env`:

   ```bash
   # Server Configuration
   PORT=8080
   HOST=localhost

   # Environment
   APP_ENV=development

   # CORS Configuration
   CORS_ORIGINS=http://localhost:3000,http://localhost:5173
   CORS_CREDENTIALS=true

   # Logging
   LOG_LEVEL=info
   LOG_FORMAT=json
   ```

### Development

Run the development server:

```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in your `.env` file).

### Building

Build the application for production:

```bash
go build -o pws-server main.go
```

### Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

### Production

Run the production server:

```bash
./pws-server
```

## Application Architecture

### Configuration Management

The application uses a singleton configuration pattern for centralized settings:

```go
// Load environment variables
config.LoadEnv()

// Setup structured logging
config.SetupLogger()

// Configure Fiber server
app := config.SetupFiber()
```

Environment variables are loaded once and accessed through thread-safe getters.

### Logging

Structured logging using Go's slog package:

```go
import "log/slog"

// Structured logging with context
slog.Info("Server starting",
    "port", config.GetPort(),
    "env", config.GetAppEnv(),
)
```

Log levels: `debug`, `info`, `warn`, `error`

### Response Handling

Standardized response structure using builder pattern:

```go
// Success response with data
return response.Success(c).
    Message("Operation completed successfully").
    Data(userData).
    Send()

// Error response
return response.Error(c).
    Message("Validation failed").
    Code("VALIDATION_ERROR").
    Details(validationErrors).
    Send()

// Paginated response
return response.Success(c).
    Message("Users retrieved successfully").
    Data(users).
    Pagination(1, 10, 100).
    Send()
```

### Middleware Stack

The application supports extensible middleware:

```go
// CORS middleware
app.Use(middleware.CORS())

// Custom middleware example
app.Use(func(c *fiber.Ctx) error {
    // Middleware logic
    return c.Next()
})
```

### Routing

Routes are organized by feature or resource:

```go
// API v1 routes
api := app.Group("/api/v1")

// Resource routes
users := api.Group("/users")
users.Get("/", handlers.GetUsers)
users.Post("/", handlers.CreateUser)
users.Get("/:id", handlers.GetUser)
users.Put("/:id", handlers.UpdateUser)
users.Delete("/:id", handlers.DeleteUser)
```

## API Reference

### Response Format

All API responses follow a consistent structure:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data
  },
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "pages": 10
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Response Format

```json
{
  "success": false,
  "message": "Error description",
  "error": {
    "code": "ERROR_CODE",
    "message": "Detailed error message",
    "details": {
      // Additional error context
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Validation Error
- `500` - Internal Server Error

## Development Guidelines

### Code Structure

1. **Package Organization**: Organize code by feature, not by layer
2. **Error Handling**: Use structured error responses
3. **Logging**: Include context in all log messages
4. **Testing**: Write tests for all business logic
5. **Documentation**: Document all public functions and types

### Adding New Routes

1. Define handler functions:

   ```go
   func GetUsers(c *fiber.Ctx) error {
       // Implementation
       return response.Success(c).
           Message("Users retrieved successfully").
           Data(users).
           Send()
   }
   ```

2. Register routes in router:

   ```go
   users := api.Group("/users")
   users.Get("/", handlers.GetUsers)
   ```

3. Add tests:
   ```go
   func TestGetUsers(t *testing.T) {
       // Test implementation
   }
   ```

### Adding Middleware

1. Create middleware function:

   ```go
   func CustomMiddleware() fiber.Handler {
       return func(c *fiber.Ctx) error {
           // Middleware logic
           return c.Next()
       }
   }
   ```

2. Register middleware:
   ```go
   app.Use(CustomMiddleware())
   ```

### Environment Variables

Add new environment variables to:

1. `.env.example` file
2. `config/env.go` for loading
3. Documentation

### Performance Considerations

1. **Connection Pooling**: Use connection pools for databases
2. **Caching**: Implement caching where appropriate
3. **Rate Limiting**: Add rate limiting for public APIs
4. **Compression**: Enable response compression
5. **Monitoring**: Add health checks and metrics

## Deployment

### Docker

Build and run with Docker:

```bash
# Build image
docker build -t pws-server .

# Run container
docker run -p 8080:8080 --env-file .env pws-server
```

### Environment Variables for Production

```bash
# Server
PORT=8080
HOST=0.0.0.0

# Environment
APP_ENV=production

# CORS
CORS_ORIGINS=https://yourdomain.com
CORS_CREDENTIALS=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Database (when implemented)
DATABASE_URL=postgresql://user:pass@localhost/pws
```

### Health Checks

The server provides health check endpoints:

- `GET /health` - Basic health check

## Dependencies

### Production Dependencies

- `github.com/gofiber/fiber/v3` - Web framework
- `github.com/joho/godotenv` - Environment variables
- `github.com/goccy/go-json` - High-performance JSON

### Development Dependencies

- Standard Go testing framework
- Go standard library packages

## Contributing

1. Follow Go coding standards and conventions
2. Write comprehensive tests for new features
3. Document all public APIs
4. Use structured logging throughout
5. Follow the established project structure
6. Ensure proper error handling

## License

This project is part of the PWS application suite.
