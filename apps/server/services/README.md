# Services Package

This package contains business logic services that handle core functionality like authentication, caching, and database operations.

## What It Does

- Handles user authentication (login, register, tokens)
- Manages password hashing and verification
- Provides caching with Redis
- Handles secure cookies for authentication
- Manages database connections and queries

## Main Services

### AuthService
Handles all authentication operations.

**Main Functions:**
- `Login(authRequest)` - Authenticates user with email/password
- `Register(registerRequest)` - Creates new user account
- `GenerateAccessToken(user)` - Creates JWT access token
- `GenerateRefreshToken(user)` - Creates JWT refresh token
- `RefreshToken(token)` - Gets new tokens using refresh token
- `GetUserByID(id)` - Retrieves user by ID
- `HashPassword(password)` - Hashes password securely
- `VerifyPassword(password, hash)` - Checks if password matches hash

**How to use:**
```go
logger := config.SetupLogger()
authService := &services.AuthService{Logger: logger}

// Login user
user, err := authService.Login(&types.AuthRequest{
    Email:    "user@example.com",
    Password: "userpassword",
})

// Register new user
user, err := authService.Register(&types.RegisterRequest{
    Username: "newuser",
    Email:    "new@example.com",
    Password: "newpassword",
})

// Generate tokens
accessToken, err := authService.GenerateAccessToken(user)
refreshToken, err := authService.GenerateRefreshToken(user)
```

### CacheService
Handles Redis caching operations.

**Main Functions:**
- `Set(key, value, expiration)` - Store data in cache
- `Get(key)` - Retrieve data from cache
- `Delete(key)` - Remove data from cache
- `Ping()` - Test Redis connection

**How to use:**
```go
cacheService := &services.CacheService{}

// Store in cache
err := cacheService.Set("user:123", userData, time.Hour)

// Get from cache
data, err := cacheService.Get("user:123")

// Delete from cache
err := cacheService.Delete("user:123")
```

### CookieService
Manages secure HTTP cookies for authentication.

**Main Functions:**
- `SetAuthCookies(c, accessToken, refreshToken)` - Set login cookies
- `ClearAuthCookies(c)` - Remove login cookies
- `GetAccessToken(c)` - Get access token from cookie
- `GetRefreshToken(c)` - Get refresh token from cookie

**How to use:**
```go
cookieService := &services.CookieService{}

// Set cookies after login
cookieService.SetAuthCookies(c, accessToken, refreshToken)

// Clear cookies on logout
cookieService.ClearAuthCookies(c)

// Get tokens from cookies
accessToken := cookieService.GetAccessToken(c)
refreshToken := cookieService.GetRefreshToken(c)
```

## Database Functions

These functions work directly with the database:

- `Ping()` - Test database connection
- `CloseDatabase()` - Close database connection
- `CloseRedisConnection()` - Close Redis connection

**How to use:**
```go
// Test connections
err := services.Ping()
if err != nil {
    log.Fatal("Database connection failed")
}

// Close connections (usually in main.go)
defer services.CloseDatabase()
defer services.CloseRedisConnection()
```

## Password Security

The AuthService uses Argon2 for secure password hashing:

```go
// Hash a password
hashedPassword, err := authService.HashPassword("userpassword")

// Verify a password
isValid, err := authService.VerifyPassword("userpassword", hashedPassword)
```

## JWT Tokens

Access tokens expire quickly (15 minutes), refresh tokens last longer (7 days):

```go
// Generate tokens
accessToken, err := authService.GenerateAccessToken(user)   // Short-lived
refreshToken, err := authService.GenerateRefreshToken(user) // Long-lived

// Refresh expired access token
newTokens, err := authService.RefreshToken(oldRefreshToken)
```

## Error Handling

All services return Go errors. Common patterns:

```go
user, err := authService.Login(request)
if err != nil {
    if errors.Is(err, lib.ErrInvalidCredentials) {
        // Handle wrong password
    } else {
        // Handle other errors
    }
}
```

## Security Features

- Passwords are hashed with Argon2
- JWT tokens are signed with secrets
- Cookies are HTTP-only and secure
- Refresh token rotation prevents replay attacks
- Database queries use parameterized statements
