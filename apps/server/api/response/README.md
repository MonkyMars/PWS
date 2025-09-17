# Response Package

The `response` package provides standardized HTTP response structures and utilities for the PWS API. This package implements a consistent response format across all API endpoints, including success responses, error handling, pagination support, and response building utilities.

## Overview

This package follows REST API best practices by providing:

- Consistent response structure across all endpoints
- Detailed error information with error codes and context
- Pagination metadata for list endpoints
- Fluent response builder pattern for easy response construction
- Standardized error codes for common scenarios

## Files

- `types.go` - Core response structures and builder pattern implementation
- `success.go` - Success response helper functions
- `error.go` - Error response helper functions
- `lib.go` - Utility functions for pagination, validation, and request handling

## Core Types

### Response Structure

The standard response structure used across all API endpoints:

```go
type Response struct {
    Success   bool       `json:"success"`               // Whether the request was successful
    Message   string     `json:"message"`               // Human-readable description
    Data      any        `json:"data,omitempty"`        // Response payload
    Error     *ErrorInfo `json:"error,omitempty"`       // Error details
    Meta      *Meta      `json:"meta,omitempty"`        // Additional metadata
    Timestamp time.Time  `json:"timestamp"`             // Response timestamp
}
```

### Error Information

```go
type ErrorInfo struct {
    Code    string         `json:"code"`               // Machine-readable error code
    Message string         `json:"message"`            // Human-readable description
    Details map[string]any `json:"details,omitempty"`  // Additional error context
    Field   string         `json:"field,omitempty"`    // Field causing error
}
```

### Pagination Metadata

```go
type Meta struct {
    Page       int  `json:"page,omitempty"`        // Current page (1-based)
    Limit      int  `json:"limit,omitempty"`       // Items per page
    Total      int  `json:"total,omitempty"`       // Total items
    TotalPages int  `json:"total_pages,omitempty"` // Total pages
    HasNext    bool `json:"has_next,omitempty"`    // More pages available
    HasPrev    bool `json:"has_prev,omitempty"`    // Previous pages available
}
```

## Builder Pattern

The package uses a fluent builder pattern for constructing responses:

```go
return NewResponse().
    Success("Request successful").
    WithData(userData).
    Send(c, fiber.StatusOK)
```

## Success Responses (`success.go`)

Functions for successful HTTP responses:

- `OK(c, data)` - 200 response with data
- `OKWithMessage(c, message, data)` - 200 response with custom message
- `Created(c, data)` - 201 response for resource creation
- `Accepted(c, message)` - 202 response for accepted requests
- `NoContent(c)` - 204 response with no content
- `Paginated(c, items, page, limit, total)` - Paginated response with metadata
- `Message(c, message)` - Simple success message without data

## Error Responses (`error.go`)

Functions for error HTTP responses:

- `BadRequest(c, message)` - 400 Bad Request
- `Unauthorized(c, message)` - 401 Unauthorized
- `Forbidden(c, message)` - 403 Forbidden
- `NotFound(c, message)` - 404 Not Found
- `Conflict(c, message)` - 409 Conflict
- `ValidationError(c, errors)` - 422 Validation Error
- `TooManyRequests(c, message)` - 429 Too Many Requests
- `InternalServerError(c, message)` - 500 Internal Server Error
- `ServiceUnavailable(c, message)` - 503 Service Unavailable

## Utility Functions (`lib.go`)

Helper functions for common operations:

### Pagination

- `ParsePaginationParams(c)` - Extract page and limit from query parameters
- `CalculateOffset(page, limit)` - Calculate database offset
- `QuickPaginate(c, items, totalCount)` - Simple pagination helper

### Validation

- `ValidateJSON(c, v)` - Validate and bind JSON request body
- `BuildValidationErrors(errors)` - Convert error map to validation error slice
- `SendValidationErrors(c, errors)` - Quick validation error response

### Context Utilities

- `GetUserID(c)` - Extract user ID from request context
- `GetUserRole(c)` - Extract user role from request context

### Response Helpers

- `ErrorWithContext(c, statusCode, message, context)` - Error response with request context
- `ConvertToInterfaceSlice(slice)` - Convert typed slices to interface slice

## Error Codes

Predefined error codes for consistent error handling:

- `VALIDATION_ERROR` - Input validation failures
- `NOT_FOUND` - Resource not found
- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Access denied
- `CONFLICT` - Resource conflict
- `INTERNAL_ERROR` - Server errors
- `BAD_REQUEST` - Invalid request format
- `TOO_MANY_REQUESTS` - Rate limiting
- `SERVICE_UNAVAILABLE` - Service temporarily down

## Usage Examples

### Simple Success Response

```go
func GetUser(c fiber.Ctx) error {
    user := getUserData()
    return response.OK(c, user)
}
```

### Error Response

```go
func CreateUser(c fiber.Ctx) error {
    if userExists {
        return response.Conflict(c, "User already exists")
    }
    // ... create user
}
```

### Paginated Response

```go
func ListUsers(c fiber.Ctx) error {
    page, limit, err := response.ParsePaginationParams(c)
    if err != nil {
        return response.BadRequest(c, err.Error())
    }
    
    users, total := getUserList(page, limit)
    return response.Paginated(c, users, page, limit, total)
}
```

### Validation Error

```go
func UpdateUser(c fiber.Ctx) error {
    var req UserRequest
    if err := response.ValidateJSON(c, &req); err != nil {
        return response.BadRequest(c, "Invalid JSON")
    }
    
    if errors := validateUser(req); len(errors) > 0 {
        return response.SendValidationErrors(c, errors)
    }
    // ... update user
}
```

### Builder Pattern Example

```go
func ComplexResponse(c fiber.Ctx) error {
    return response.NewResponse().
        Success("Operation completed successfully").
        WithData(data).
        WithMeta(&response.Meta{
            Page:  1,
            Limit: 10,
            Total: 100,
        }).
        Send(c, fiber.StatusOK)
}
```

## Response Examples

### Success Response

```json
{
  "success": true,
  "message": "Request successful",
  "data": {"id": 1, "name": "John"},
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Response

```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The given data was invalid",
    "details": {
      "validation_errors": [
        {"field": "email", "message": "Invalid email format"}
      ]
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Paginated Response

```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {
    "items": [{"id": 1}, {"id": 2}],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false
    }
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Best Practices

1. **Consistent Error Handling**: Always use the package functions for error responses
2. **Meaningful Messages**: Provide clear, user-friendly error messages
3. **Validation Errors**: Use structured validation errors for form validation
4. **Pagination**: Use the pagination utilities for consistent list responses
5. **Builder Pattern**: Use the response builder for complex responses
6. **Error Codes**: Use the predefined error codes for machine-readable error handling

## Dependencies

- `github.com/gofiber/fiber/v3` - Web framework for HTTP handling
- `github.com/MonkyMars/PWS/types` - Shared type definitions
- Standard library packages for JSON and time handling