# Middleware Package

The `middleware` package provides HTTP middleware functions for the PWS application. This package contains reusable middleware components for cross-cutting concerns such as CORS handling, authentication, logging, rate limiting, and request validation.

## Overview

Middleware functions in this package follow the Fiber v3 middleware pattern and can be easily integrated into the application's HTTP request processing pipeline. Each middleware function is designed to be modular, testable, and configurable.

## Current Implementation

### CORS Middleware (`cors.go`)

Handles Cross-Origin Resource Sharing configuration to allow or restrict web browsers from making requests to the server from different origins.

#### Features

- **Configurable Origins**: Set allowed origins for cross-origin requests
- **HTTP Methods**: Configure which HTTP methods are allowed
- **Headers Support**: Specify which headers can be used in requests
- **Credentials Handling**: Support for cookies and authorization headers
- **Preflight Requests**: Automatic handling of OPTIONS preflight requests

#### Configuration

The CORS middleware is configured with development-friendly defaults:

```go
func SetupCORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: true,
    })
}
```

#### Usage

```go
app := fiber.New()
app.Use(middleware.SetupCORS())
```

## Middleware Pattern

All middleware functions in this package follow a consistent pattern:

```go
func MiddlewareName(config ...Config) fiber.Handler {
    // Set default configuration
    cfg := configDefault(config...)
    
    return func(c fiber.Ctx) error {
        // Middleware logic here
        
        // Continue to next handler
        return c.Next()
    }
}
```

## Planned Middleware

Future middleware components may include:

- **Authentication**: JWT token validation
- **Rate Limiting**: Request rate limiting per IP or user
- **Logging**: Request/response logging with structured data
- **Compression**: Response compression (gzip, deflate)
- **Security Headers**: Security-related HTTP headers
- **Request ID**: Unique request identifier for tracing

## Best Practices

1. **Stateless**: Middleware should be stateless when possible
2. **Configuration**: Provide configurable options with sensible defaults
3. **Error Handling**: Handle errors gracefully and continue or abort appropriately
4. **Performance**: Minimize overhead in middleware execution
5. **Testing**: Write unit tests for middleware logic
6. **Documentation**: Document configuration options and behavior

## Dependencies

- `github.com/gofiber/fiber/v3` - Web framework
- `github.com/gofiber/fiber/v3/middleware/cors` - CORS middleware implementation