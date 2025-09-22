# Middleware Package

This package contains HTTP middleware functions that process requests before they reach the route handlers. Middleware can modify requests, add data to the context, or block requests entirely.

## What It Does

- Handles cross-origin resource sharing (CORS)
- Validates authentication tokens
- Sets user information in request context
- Blocks unauthorized access to protected routes
- Provides security and request processing layers

## Main Files

### `cors.go`
Sets up CORS (Cross-Origin Resource Sharing) to allow web browsers to make requests from different domains.

**Functions:**

**`SetupCORS()`** - Returns CORS middleware
```go
// Sets up CORS headers for browser requests
// Allows the frontend to make API calls from a different port/domain
func SetupCORS() fiber.Handler
```

**How to use:**
```go
app.Use(middleware.SetupCORS())
```

**What it does:**
- Allows requests from frontend domains
- Sets proper CORS headers
- Handles preflight OPTIONS requests
- Enables credential sharing (cookies)

### `auth.go`
Handles authentication for protected routes.

**Functions:**

**`AuthMiddleware()`** - Returns authentication middleware
```go
// Validates access tokens and sets user info in context
// Blocks requests with missing or invalid tokens
func AuthMiddleware() fiber.Handler
```

**How to use:**
```go
// Apply to specific route
app.Get("/protected", middleware.AuthMiddleware(), handler)

// Apply to route group
protected := app.Group("/api", middleware.AuthMiddleware())
```

## How Middleware Works

Middleware functions run **before** your route handlers. They can:

1. **Modify the request** (add headers, parse data)
2. **Add data to context** (user info, request ID)
3. **Block the request** (authentication failure)
4. **Continue to next middleware** (everything is OK)

### Middleware Flow

```
Request -> CORS Middleware -> Auth Middleware -> Route Handler -> Response
```

## Authentication Middleware Details

The auth middleware performs these steps:

1. **Extract token** from HTTP-only cookie
2. **Validate token** (signature, expiration)
3. **Check blacklist** (logout, security)
4. **Set user info** in request context
5. **Continue** to route handler or **block** request

### Token Validation Process

```go
// 1. Get token from cookie
token := c.Cookies("access_token")

// 2. Parse and validate token
claims, err := authService.ParseToken(token, true)

// 3. Check if token is blacklisted
blacklisted, err := cacheService.IsTokenBlacklisted(claims.Jti)

// 4. Set user info in context
c.Locals("claims", claims)

// 5. Continue to handler
return c.Next()
```

## Using Protected Routes

When a route uses `AuthMiddleware()`, the handler can access user information:

```go
func ProtectedHandler(c fiber.Ctx) error {
    // Get user claims from middleware
    claimsInterface := c.Locals("claims")
    claims, ok := claimsInterface.(*types.AuthClaims)
    if !ok {
        return response.Unauthorized(c, "Unauthorized")
    }
    
    // Use user information
    userID := claims.Sub
    userEmail := claims.Email
    
    // Process request for this user
    return response.Success(c, data)
}
```

## CORS Configuration

CORS allows the frontend (running on a different port) to make requests to the API:

```go
// Frontend (localhost:5173) can call API (localhost:8080)
// Without CORS, browsers would block these requests
app.Use(middleware.SetupCORS())
```

**CORS headers set:**
- `Access-Control-Allow-Origin` - Which domains can make requests
- `Access-Control-Allow-Methods` - Which HTTP methods are allowed
- `Access-Control-Allow-Headers` - Which headers are allowed
- `Access-Control-Allow-Credentials` - Whether cookies are allowed

## Security Features

### Token Blacklisting
- Tokens are checked against a blacklist in Redis
- Logged out tokens are immediately blocked
- Protects against token reuse attacks

### Graceful Degradation
- If Redis is down, requests are blocked for security
- Logs security events for monitoring
- Fails closed rather than open

### Attack Detection
- Logs attempts to use blacklisted tokens
- Includes client IP and user agent
- Helps detect token reuse attacks

## Error Handling

The auth middleware returns appropriate HTTP status codes:

```go
// Missing token
return response.Unauthorized(c, "Missing access token")

// Invalid token
return response.Unauthorized(c, "Invalid or expired access token")

// Blacklisted token
return response.Unauthorized(c, "Token has been revoked")

// Redis error
return response.InternalServerError(c, "Authentication service temporarily unavailable")
```

## Middleware Order

Middleware order matters. Apply them in this sequence:

```go
// 1. CORS first (for browser compatibility)
app.Use(middleware.SetupCORS())

// 2. Logging middleware
app.Use(logger.HTTPMiddleware())

// 3. Auth middleware on protected routes only
protected := app.Group("/api", middleware.AuthMiddleware())
```

## Creating Custom Middleware

Follow this pattern to create new middleware:

```go
func CustomMiddleware() fiber.Handler {
    return func(c fiber.Ctx) error {
        // Do something before the handler
        
        // Optionally modify request
        c.Set("X-Custom-Header", "value")
        
        // Continue to next middleware/handler
        if err := c.Next(); err != nil {
            return err
        }
        
        // Do something after the handler
        return nil
    }
}
```

## Best Practices

1. **Order matters** - Apply middleware in logical order
2. **Fail securely** - Block requests when security checks fail
3. **Log security events** - Track authentication failures
4. **Handle errors gracefully** - Don't crash on middleware errors
5. **Keep middleware focused** - Each middleware should have one responsibility
6. **Use context** - Pass data between middleware and handlers via `c.Locals()`

## Common Use Cases

**Public routes** (no auth required):
```go
app.Post("/auth/login", internal.Login)
app.Post("/auth/register", internal.Register)
```

**Protected routes** (auth required):
```go
app.Get("/auth/me", middleware.AuthMiddleware(), internal.Me)
app.Post("/users", middleware.AuthMiddleware(), internal.CreateUser)
```

**Route groups** (apply middleware to multiple routes):
```go
api := app.Group("/api", middleware.AuthMiddleware())
api.Get("/users", internal.GetUsers)
api.Post("/users", internal.CreateUser)
```
