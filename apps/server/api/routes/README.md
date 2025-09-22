# Routes Package

This package defines all HTTP route endpoints for the API. It organizes routes into logical groups and connects them to their corresponding handler functions.

## What It Does

- Defines all API endpoints and their HTTP methods
- Groups related routes together (auth, app, etc.)
- Connects routes to handler functions
- Sets up middleware for specific routes
- Provides a central place to see all available endpoints

## Main Files

### `auth.go`
Contains all authentication-related routes.

**Routes:**
- `POST /auth/login` - User login
- `POST /auth/register` - User registration  
- `POST /auth/refresh` - Refresh access token
- `POST /auth/logout` - User logout
- `GET /auth/me` - Get current user info (protected)

**How it's structured:**
```go
func SetupAuthRoutes(app *fiber.App) {
    // Create auth route group
    auth := app.Group("/auth")
    
    // Public routes
    auth.Post("/login", internal.Login)
    auth.Post("/register", internal.Register)
    auth.Post("/refresh", internal.RefreshToken)
    auth.Post("/logout", internal.Logout)
    
    // Protected routes (requires authentication)
    auth.Get("/me", middleware.AuthMiddleware(), internal.Me)
}
```

### `app.go`
Contains general application routes like health checks.

**Routes (development only):**
- `GET /health` - Server health check
- `GET /health/database` - Database health check
- `/*` - Catch-all for 404 errors

**How it's structured:**
```go
func SetupAppRoutes(app *fiber.App) {
    cfg := config.Get()
    if cfg.Environment == "development" {
        app.Get("/health", internal.GetSystemHealth)
        app.Get("/health/database", internal.GetDatabaseHealth)
    }
    app.Use(internal.NotFoundHandler)
}
```

## How Routes Work

Routes are set up in the main API router (`api/router.go`):

```go
func SetupRoutes(app *fiber.App, logger *config.Logger) {
    // Authentication routes
    routes.SetupAuthRoutes(app)
    
    // Health check and fallback routes (must be last)
    routes.SetupAppRoutes(app)
}
```

## Route Groups

**Auth Group (`/auth`):**
All authentication endpoints start with `/auth/`. This makes it easy to apply auth-specific middleware or rate limiting.

**Health Routes:**
Only available in development mode for security. In production, use external monitoring instead.

## Middleware

Routes can have middleware applied:

```go
// Apply middleware to specific route
auth.Get("/me", middleware.AuthMiddleware(), internal.Me)

// Apply middleware to entire group
auth.Use(middleware.AuthMiddleware())
```

## Adding New Routes

1. **Create route function in appropriate file:**
```go
func SetupUserRoutes(app *fiber.App) {
    users := app.Group("/users")
    
    users.Get("/", internal.GetUsers)
    users.Get("/:id", internal.GetUser)
    users.Post("/", middleware.AuthMiddleware(), internal.CreateUser)
    users.Put("/:id", middleware.AuthMiddleware(), internal.UpdateUser)
    users.Delete("/:id", middleware.AuthMiddleware(), internal.DeleteUser)
}
```

2. **Add to main router:**
```go
func SetupRoutes(app *fiber.App, logger *config.Logger) {
    routes.SetupAuthRoutes(app)
    routes.SetupUserRoutes(app)  // Add new routes
    routes.SetupAppRoutes(app)   // Keep this last
}
```

3. **Create handlers in `api/internal/`:**
```go
func GetUsers(c fiber.Ctx) error {
    // Handle GET /users
}

func CreateUser(c fiber.Ctx) error {
    // Handle POST /users
}
```

## Route Parameters

**Path parameters:**
```go
app.Get("/users/:id", handler)  // /users/123
userID := c.Params("id")
```

**Query parameters:**
```go
app.Get("/users", handler)      // /users?page=1&limit=10
page := c.Query("page")
```

**Request body:**
```go
var req types.CreateUserRequest
if err := c.Bind().Body(&req); err != nil {
    return response.BadRequest(c, "Invalid request body")
}
```

## Protected Routes

Routes that require authentication use the `AuthMiddleware`:

```go
auth.Get("/me", middleware.AuthMiddleware(), internal.Me)
```

The middleware:
- Checks for valid access token in cookies
- Validates the JWT token
- Sets user information in the request context
- Returns 401 if authentication fails

## Error Handling

Route handlers should return appropriate responses:

```go
func GetUser(c fiber.Ctx) error {
    userID := c.Params("id")
    if userID == "" {
        return response.BadRequest(c, "User ID is required")
    }
    
    user, err := userService.GetByID(userID)
    if err != nil {
        return response.NotFound(c, "User not found")
    }
    
    return response.Success(c, user)
}
```

## Security Considerations

1. **Auth routes:** Public endpoints are limited to prevent abuse
2. **Health routes:** Only enabled in development
3. **Protected routes:** Always use `AuthMiddleware()` for sensitive operations
4. **Input validation:** Validate all parameters and request bodies
5. **Rate limiting:** Consider adding rate limiting middleware for public endpoints

## Route Order

Route order matters in Fiber. More specific routes should come before general ones:

```go
app.Get("/auth/me", handler)     // Specific route first
app.Get("/auth/*", handler)      // Wildcard route after
```

The 404 handler (`app.Use(internal.NotFoundHandler)`) must be last to catch all unmatched routes.