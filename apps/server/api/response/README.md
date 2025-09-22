# Response Package

This package provides standardized API response formatting for all HTTP endpoints. It ensures consistent response structure across the entire application.

## What It Does

- Creates consistent JSON responses for success and error cases
- Handles pagination for list endpoints
- Provides validation error formatting
- Manages HTTP status codes automatically
- Includes timestamps and metadata

## Response Structure

All API responses follow this format:

```json
{
  "success": true,
  "message": "Request successful",
  "data": { },
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "pages": 10
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Main Functions

### Success Responses

**`Success(c, data)`** - Basic success response
```go
return response.Success(c, user)
```

**`SuccessWithMessage(c, message, data)`** - Success with custom message
```go
return response.SuccessWithMessage(c, "User updated successfully", user)
```

**`Created(c, data)`** - For newly created resources (201 status)
```go
return response.Created(c, newUser)
```

**`Paginated(c, items, page, limit, total)`** - For paginated lists
```go
return response.Paginated(c, users, 1, 10, 100)
```

### Error Responses

**`BadRequest(c, message)`** - 400 Bad Request
```go
return response.BadRequest(c, "Invalid input data")
```

**`Unauthorized(c, message)`** - 401 Unauthorized
```go
return response.Unauthorized(c, "Invalid credentials")
```

**`Forbidden(c, message)`** - 403 Forbidden
```go
return response.Forbidden(c, "Access denied")
```

**`NotFound(c, message)`** - 404 Not Found
```go
return response.NotFound(c, "User not found")
```

**`Conflict(c, message)`** - 409 Conflict
```go
return response.Conflict(c, "Email already exists")
```

**`InternalServerError(c, message)`** - 500 Internal Server Error
```go
return response.InternalServerError(c, "Something went wrong")
```

**`ServiceUnavailable(c, message)`** - 503 Service Unavailable
```go
return response.ServiceUnavailable(c, "Database temporarily unavailable")
```

### Validation Errors

**`SendValidationError(c, errors)`** - 422 Unprocessable Entity
```go
validationErrors := []types.ValidationError{
    {
        Field:   "email",
        Message: "Email is required",
        Value:   "",
    },
}
return response.SendValidationError(c, validationErrors)
```

## How to Use

### In a Route Handler

```go
func GetUser(c fiber.Ctx) error {
    userID := c.Params("id")
    
    // Validate input
    if userID == "" {
        return response.BadRequest(c, "User ID is required")
    }
    
    // Get user from database
    user, err := userService.GetByID(userID)
    if err != nil {
        if errors.Is(err, lib.ErrNotFound) {
            return response.NotFound(c, "User not found")
        }
        return response.InternalServerError(c, "Failed to retrieve user")
    }
    
    // Return success response
    return response.Success(c, user)
}
```

### With Pagination

```go
func GetUsers(c fiber.Ctx) error {
    // Parse pagination parameters
    page, limit, err := response.ParsePaginationParams(c)
    if err != nil {
        return response.BadRequest(c, err.Error())
    }
    
    // Get users from database
    users, total, err := userService.GetPaginated(page, limit)
    if err != nil {
        return response.InternalServerError(c, "Failed to retrieve users")
    }
    
    // Return paginated response
    return response.Paginated(c, users, page, limit, total)
}
```

### With Validation

```go
func CreateUser(c fiber.Ctx) error {
    var req types.CreateUserRequest
    if err := c.Bind().Body(&req); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    // Validate fields
    var validationErrors []types.ValidationError
    if req.Email == "" {
        validationErrors = append(validationErrors, types.ValidationError{
            Field:   "email",
            Message: "Email is required",
            Value:   req.Email,
        })
    }
    if req.Username == "" {
        validationErrors = append(validationErrors, types.ValidationError{
            Field:   "username",
            Message: "Username is required",
            Value:   req.Username,
        })
    }
    
    if len(validationErrors) > 0 {
        return response.SendValidationError(c, validationErrors)
    }
    
    // Create user
    user, err := userService.Create(&req)
    if err != nil {
        return response.InternalServerError(c, "Failed to create user")
    }
    
    return response.Created(c, user)
}
```

## Helper Functions

**`ParsePaginationParams(c)`** - Extract page and limit from query params
```go
page, limit, err := response.ParsePaginationParams(c)
// Default: page=1, limit=10, max limit=100
```

**`CalculateOffset(page, limit)`** - Convert page number to database offset
```go
offset := response.CalculateOffset(page, limit)
```

**`GetUserID(c)`** - Extract authenticated user ID from context
```go
userID, err := response.GetUserID(c)
```

## Builder Pattern

For complex responses, use the response builder:

```go
return NewResponse().
    Success("User retrieved successfully").
    WithData(user).
    WithMeta(&types.Meta{
        RequestID: "req-123",
        Version:   "v1",
    }).
    Send(c, fiber.StatusOK)
```

## Best Practices

1. Always use the response functions instead of raw JSON
2. Provide clear, user-friendly error messages
3. Use appropriate HTTP status codes
4. Include validation details for client debugging
5. Keep response data structure consistent
6. Use pagination for lists that could grow large