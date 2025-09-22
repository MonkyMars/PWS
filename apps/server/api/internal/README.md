# Internal Package

This package contains the actual handler functions that process HTTP requests. These are the functions that do the real work when someone calls an API endpoint.

## What It Does

- Contains handler functions for all API endpoints
- Processes incoming HTTP requests
- Validates request data
- Calls services to perform business logic
- Returns formatted responses using the response package
- Handles authentication and authorization

## Main Files

### `auth.go`
Contains authentication handler functions.

**Functions:**

**`Login(c fiber.Ctx)`** - Handles user login
```go
// POST /auth/login
// Expects: {"email": "user@example.com", "password": "password"}
// Returns: User data and sets authentication cookies
func Login(c fiber.Ctx) error
```

**`Register(c fiber.Ctx)`** - Handles user registration
```go
// POST /auth/register  
// Expects: {"username": "user", "email": "user@example.com", "password": "password"}
// Returns: New user data and sets authentication cookies
func Register(c fiber.Ctx) error
```

**`RefreshToken(c fiber.Ctx)`** - Refreshes expired access tokens
```go
// POST /auth/refresh
// Uses refresh token from cookies to get new access token
// Returns: New tokens and rotates refresh token
func RefreshToken(c fiber.Ctx) error
```

**`Logout(c fiber.Ctx)`** - Handles user logout
```go
// POST /auth/logout
// Blacklists current tokens and clears cookies
// Returns: Success message
func Logout(c fiber.Ctx) error
```

**`Me(c fiber.Ctx)`** - Gets current user info (protected route)
```go
// GET /auth/me
// Requires: Valid access token
// Returns: Current user information
func Me(c fiber.Ctx) error
```

### `app.go`
Contains general application handler functions.

**Functions:**

**`GetSystemHealth(c fiber.Ctx)`** - Returns server health status
```go
// GET /health (development only)
// Returns: Server status, uptime, memory usage, database status
func GetSystemHealth(c fiber.Ctx) error
```

**`GetDatabaseHealth(c fiber.Ctx)`** - Checks database connection
```go
// GET /health/database (development only)
// Returns: Database connection status and response time
func GetDatabaseHealth(c fiber.Ctx) error
```

**`NotFoundHandler(c fiber.Ctx)`** - Handles 404 errors
```go
// Catches all unmatched routes
// Returns: 404 error with standard format
func NotFoundHandler(c fiber.Ctx) error
```

## How Handlers Work

All handler functions follow the same pattern:

1. **Extract request data** (parameters, body, etc.)
2. **Validate the input** (required fields, format, etc.)
3. **Call services** to perform business logic
4. **Handle errors** appropriately
5. **Return response** using response package

### Example Handler Pattern

```go
func CreateUser(c fiber.Ctx) error {
    // 1. Extract request data
    var req types.CreateUserRequest
    if err := c.Bind().Body(&req); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    // 2. Validate input
    if req.Email == "" {
        return response.BadRequest(c, "Email is required")
    }
    
    // 3. Call service
    authService := &services.AuthService{}
    user, err := authService.CreateUser(&req)
    if err != nil {
        return response.InternalServerError(c, "Failed to create user")
    }
    
    // 4. Return response
    return response.Created(c, user)
}
```

## Authentication Flow

### Login Process
1. User sends email/password
2. `Login` handler validates credentials
3. If valid, generates access and refresh tokens
4. Sets secure HTTP-only cookies
5. Returns user data

### Protected Route Access
1. Client makes request with cookies
2. `AuthMiddleware` validates access token
3. If valid, sets user info in context
4. Handler accesses user info via `c.Locals("claims")`
5. Processes request for authenticated user

### Token Refresh
1. When access token expires, client calls `/auth/refresh`
2. `RefreshToken` handler validates refresh token
3. Generates new access token and rotates refresh token
4. Sets new cookies
5. Client can continue with new tokens

## Input Validation

Handlers validate input at multiple levels:

```go
// Required field validation
if strings.TrimSpace(authRequest.Email) == "" {
    return response.SendValidationError(c, []types.ValidationError{
        {
            Field:   "email",
            Message: "Email is required",
            Value:   authRequest.Email,
        },
    })
}

// Format validation
if len(registerRequest.Password) < 6 {
    return response.SendValidationError(c, []types.ValidationError{
        {
            Field:   "password", 
            Message: "Password must be at least 6 characters long",
            Value:   registerRequest.Password,
        },
    })
}
```

## Error Handling

Handlers use the response package for consistent error handling:

```go
// Bad request (400)
return response.BadRequest(c, "Invalid input data")

// Unauthorized (401) 
return response.Unauthorized(c, "Invalid credentials")

// Not found (404)
return response.NotFound(c, "User not found")

// Conflict (409)
return response.Conflict(c, "Email already exists")

// Internal server error (500)
return response.InternalServerError(c, "Something went wrong")

// Validation errors (422)
return response.SendValidationError(c, validationErrors)
```

## Working with Services

Handlers create service instances and call their methods:

```go
func Login(c fiber.Ctx) error {
    // Create service instance
    authService := &services.AuthService{}
    cookieService := &services.CookieService{}
    
    // Call service methods
    user, err := authService.Login(&authRequest)
    if err != nil {
        return response.Unauthorized(c, "Invalid credentials")
    }
    
    // Generate tokens
    accessToken, err := authService.GenerateAccessToken(user)
    refreshToken, err := authService.GenerateRefreshToken(user)
    
    // Set cookies
    cookieService.SetAuthCookies(c, accessToken, refreshToken)
    
    return response.Success(c, user)
}
```

## Security Considerations

1. **Input validation:** Always validate all input data
2. **Authentication:** Check authentication status for protected routes
3. **Error messages:** Don't leak sensitive information in error messages
4. **Logging:** Log security events for monitoring
5. **Rate limiting:** Consider implementing rate limiting for auth endpoints

## Request Context

Protected routes can access user information from the request context:

```go
func Me(c fiber.Ctx) error {
    // Get user claims from middleware
    claimsInterface := c.Locals("claims")
    claims, ok := claimsInterface.(*types.AuthClaims)
    if !ok {
        return response.Unauthorized(c, "Unauthorized")
    }
    
    // Use user ID from claims
    userID := claims.Sub
    user, err := authService.GetUserByID(userID)
    
    return response.Success(c, user)
}
```

## Adding New Handlers

1. **Define the handler function:**
```go
func GetUsers(c fiber.Ctx) error {
    // Implementation here
}
```

2. **Add to routes package:**
```go
users.Get("/", internal.GetUsers)
```

3. **Follow the handler pattern:**
   - Extract and validate input
   - Call appropriate services
   - Handle errors properly
   - Return formatted response

## Testing Handlers

Handlers can be tested using Fiber's test utilities:

```go
func TestLogin(t *testing.T) {
    app := fiber.New()
    app.Post("/auth/login", Login)
    
    req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{
        "email": "test@example.com",
        "password": "password"
    }`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```
