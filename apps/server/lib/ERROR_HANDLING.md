# Standardized Error Handling

This document explains the new standardized error handling system implemented to improve consistency and developer experience across the PWS API.

## Overview

The standardized error handling system provides:

- **Consistent Error Responses**: All services return predictable HTTP status codes and error messages
- **Centralized Error Management**: Single point for mapping service errors to HTTP responses
- **Better Developer Experience**: Clear error patterns and easy-to-use handler functions
- **Comprehensive Logging**: Automatic error logging with request context

## Error Handler Functions

### 1. `HandleServiceError(c fiber.Ctx, err error)`

**Primary function for handling service layer errors**

```go
// Example usage in handlers
user, err := ar.authService.Login(&authRequest)
if err != nil {
    return lib.HandleServiceError(c, err)
}
```

### 2. `HandleAuthError(c fiber.Ctx, err error, context string)`

**Specialized function for authentication errors with audit logging**

```go
// Example usage for authentication scenarios
user, err := ar.authService.Login(&authRequest)
if err != nil {
    return lib.HandleAuthError(c, err, "login")
}
```

### 3. `HandleValidationError(c fiber.Ctx, err error, field string)`

**Function for request validation errors**

```go
// Example usage for validation failures
if err := validator.Validate(request); err != nil {
    return lib.HandleValidationError(c, err, "email")
}
```

## Error Mapping

The error handler automatically maps service errors to appropriate HTTP responses:

### Authentication & Authorization (401)

```go
ErrInvalidCredentials → "Invalid email or password"
ErrInvalidToken       → "Invalid or expired token"
ErrExpiredToken       → "Invalid or expired token"
ErrUnauthorized       → "Authentication required"
```

### Forbidden Access (403)

```go
ErrInsufficientPermissions → "You do not have permission to perform this action"
ErrFileAccess             → "File access denied"
```

### Not Found (404)

```go
ErrUserNotFound   → "User not found"
ErrFileNotFound   → "File not found"
ErrFolderNotFound → "Folder not found"
```

### Conflict (409)

```go
ErrUserAlreadyExists → "User with this email already exists"
ErrUsernameTaken     → "Username is already taken"
```

### Bad Request (400)

```go
ErrInvalidInput  → "Invalid input data"
ErrInvalidFormat → "Invalid input data"
ErrMissingField  → "Required field is missing"
```

### Service Unavailable (503)

```go
ErrServiceUnavailable → "Service temporarily unavailable"
```

### Internal Server Error (500)

```go
ErrTokenGeneration    → "Failed to generate authentication token"
ErrTokenRefresh       → "Failed to refresh token"
ErrPasswordHashing    → "User account creation failed"
ErrUserCreation       → "User account creation failed"
ErrFileUpload         → "File upload failed"
ErrFolderCreation     → "Folder creation failed"
ErrDatabaseConnection → "Database connection error"
ErrExternalService    → "External service error"
```

## Before and After Examples

### Before (Inconsistent)

```go
func (ar *AuthRoutes) Login(c fiber.Ctx) error {
    user, err := ar.authService.Login(&authRequest)
    if err != nil {
        if errors.Is(err, lib.ErrInvalidCredentials) {
            return response.Unauthorized(c, "Invalid email or password")
        }
        return response.InternalServerError(c, "An error occurred during login")
    }

    accessToken, err := ar.authService.GenerateAccessToken(user)
    if err != nil {
        return response.InternalServerError(c, "Failed to generate access token")
    }

    // ... rest of handler
}
```

### After (Standardized)

```go
func (ar *AuthRoutes) Login(c fiber.Ctx) error {
    user, err := ar.authService.Login(&authRequest)
    if err != nil {
        return lib.HandleAuthError(c, err, "login")
    }

    accessToken, err := ar.authService.GenerateAccessToken(user)
    if err != nil {
        return lib.HandleServiceError(c, lib.ErrTokenGeneration)
    }

    // ... rest of handler
}
```

## Benefits

### 1. **Consistency**

- All similar errors return the same HTTP status codes and messages
- Predictable API behavior for frontend developers

### 2. **Maintainability**

- Error handling logic centralized in one place
- Easy to update error messages across the entire application

### 3. **Debugging**

- Automatic request context logging
- Consistent error format for monitoring and debugging

### 4. **Developer Experience**

- Simple, clear functions to use in handlers
- Reduced boilerplate code
- Self-documenting error handling patterns

## Migration Guide

To migrate existing handlers to use standardized error handling:

1. **Replace direct response calls**:

   ```go
   // Old
   return response.Unauthorized(c, "Invalid credentials")

   // New
   return lib.HandleServiceError(c, lib.ErrInvalidCredentials)
   ```

2. **Use appropriate handler function**:

   - General service errors: `HandleServiceError()`
   - Authentication errors: `HandleAuthError()`
   - Validation errors: `HandleValidationError()`

3. **Update service layer to return standard errors**:
   ```go
   // Service layer should return standard lib errors
   if !validCredentials {
       return nil, lib.ErrInvalidCredentials
   }
   ```

## Error Logging

The error handler automatically logs errors with request context:

```go
// Automatic logging includes:
{
    "error": "invalid credentials",
    "method": "POST",
    "path": "/auth/login",
    "ip": "192.168.1.100",
    "timestamp": "2025-10-10T10:30:00Z"
}
```

This standardized error handling system significantly improves the consistency and maintainability of error responses across the PWS API while providing a better developer experience.
