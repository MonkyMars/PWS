# Types Package

This package defines all the data structures and types used throughout the application. It provides a centralized location for type definitions that are shared between different parts of the server.

## What It Does

- Defines request and response structures for API endpoints
- Provides authentication-related types (claims, tokens)
- Contains database model definitions
- Defines validation error structures
- Ensures type consistency across the entire application

## Main Files

### `auth.go`
Contains authentication-related types.

**User Authentication:**
```go
type AuthRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

**JWT Claims:**
```go
type AuthClaims struct {
    Sub   uuid.UUID `json:"sub"`     // User ID
    Email string    `json:"email"`   // User email
    Role  string    `json:"role"`    // User role
    Iat   time.Time `json:"iat"`     // Issued at
    Exp   time.Time `json:"exp"`     // Expires at
    Jti   uuid.UUID `json:"jti"`     // JWT ID (for blacklisting)
}
```

**Response Types:**
```go
type AuthResponse struct {
    User         *User  `json:"user"`
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct {
    Message string `json:"message"`
}
```

**Password Hashing:**
```go
type ArgonParams struct {
    Memory  uint32  // Memory usage in KB
    Time    uint32  // Number of iterations
    Threads uint8   // Number of threads
    KeyLen  uint32  // Length of derived key
    SaltLen uint32  // Length of salt
}
```

### `response.go`
Contains API response structures.

**Standard Response:**
```go
type APIResponse struct {
    Success   bool        `json:"success"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Error     *ErrorInfo  `json:"error,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
}
```

**Error Information:**
```go
type ErrorInfo struct {
    Code    string                 `json:"code"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

type ValidationError struct {
    Field   string      `json:"field"`
    Message string      `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}
```

**Pagination:**
```go
type PaginationMeta struct {
    Page       int  `json:"page"`
    Limit      int  `json:"limit"`
    Total      int  `json:"total"`
    TotalPages int  `json:"total_pages"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_prev"`
}
```

### `app.go`
Contains application-specific types.

**Health Check:**
```go
type HealthResponse struct {
    Status            string            `json:"status"`
    Message           string            `json:"message"`
    ApplicationUptime string            `json:"application_uptime"`
    Services          map[string]string `json:"services"`
    Metrics           HealthMetrics     `json:"metrics"`
}

type HealthMetrics struct {
    MemoryUsageMB float64 `json:"memory_usage_mb"`
    GoRoutines    int     `json:"go_routines"`
    RequestCount  int64   `json:"request_count"`
}

type DatabaseHealthResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Elapsed string `json:"elapsed"`
}
```

**User Model:**
```go
type User struct {
    Id        uuid.UUID `json:"id" db:"id"`
    Username  string    `json:"username" db:"username"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`          // Hidden from JSON
    Role      string    `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### `queries.go`
Contains database query result types.

**Query Results:**
```go
type UserQueryResult struct {
    Users []User `json:"users"`
    Total int    `json:"total"`
}
```

## How to Use Types

### In Request Handlers
```go
func Login(c fiber.Ctx) error {
    var authRequest types.AuthRequest
    if err := c.Bind().Body(&authRequest); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    // Use authRequest.Email and authRequest.Password
}
```

### In Service Functions
```go
func (s *AuthService) Login(req *types.AuthRequest) (*types.User, error) {
    // Function accepts typed request
    user := &types.User{
        Email: req.Email,
        // ... other fields
    }
    return user, nil
}
```

### In Response Building
```go
func Login(c fiber.Ctx) error {
    user, err := authService.Login(&authRequest)
    if err != nil {
        return response.Unauthorized(c, "Invalid credentials")
    }
    
    authResponse := types.AuthResponse{
        User:         user,
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }
    
    return response.Success(c, authResponse)
}
```

### Working with JWT Claims
```go
func Me(c fiber.Ctx) error {
    claimsInterface := c.Locals("claims")
    claims, ok := claimsInterface.(*types.AuthClaims)
    if !ok {
        return response.Unauthorized(c, "Invalid claims")
    }
    
    userID := claims.Sub
    userEmail := claims.Email
    // Use claims data
}
```

## JSON Tags

Types use JSON tags to control how they're serialized:

```go
type User struct {
    Id       uuid.UUID `json:"id"`         // Included in JSON
    Password string    `json:"-"`          // Excluded from JSON
    Email    string    `json:"email"`      // Included in JSON
}
```

**Common JSON tags:**
- `json:"field_name"` - Include with custom name
- `json:"-"` - Exclude from JSON entirely
- `json:"field,omitempty"` - Exclude if empty/zero value

## Database Tags

Types that map to database tables use `db` tags:

```go
type User struct {
    Id       uuid.UUID `json:"id" db:"id"`
    Username string    `json:"username" db:"username"`
    Email    string    `json:"email" db:"email"`
}
```

## Validation

Types can include validation rules in comments or use validation libraries:

```go
type RegisterRequest struct {
    Username string `json:"username"` // Required, 3-50 characters
    Email    string `json:"email"`    // Required, valid email format
    Password string `json:"password"` // Required, minimum 6 characters
}
```

## Best Practices

1. **Consistent naming** - Use clear, descriptive names
2. **JSON tags** - Always include appropriate JSON tags
3. **Password security** - Use `json:"-"` for password fields
4. **Required vs Optional** - Use pointers for optional fields in update requests
5. **Documentation** - Add comments explaining complex types
6. **Validation** - Document validation rules in comments
7. **Grouping** - Keep related types in the same file

## Creating New Types

When adding new types:

1. **Choose the right file** based on the type's purpose
2. **Follow naming conventions** (PascalCase for exported types)
3. **Add JSON tags** for all fields that should be serialized
4. **Add database tags** if the type maps to a database table
5. **Document the type** with comments
6. **Consider validation needs** and document requirements

Example new type:
```go
// CreatePostRequest represents a request to create a new blog post
type CreatePostRequest struct {
    Title   string `json:"title"`   // Required, 1-200 characters
    Content string `json:"content"` // Required, 1-10000 characters
    Tags    []string `json:"tags,omitempty"` // Optional, max 10 tags
}

// Post represents a blog post in the database
type Post struct {
    Id        uuid.UUID `json:"id" db:"id"`
    Title     string    `json:"title" db:"title"`
    Content   string    `json:"content" db:"content"`
    AuthorId  uuid.UUID `json:"author_id" db:"author_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

This package ensures type safety and consistency across the entire application, making the code more maintainable and less prone to errors.