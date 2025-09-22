# Lib Package

This package contains shared utilities, constants, and helper functions used throughout the application. It provides common functionality that doesn't belong to any specific domain.

## What It Does

- Defines application-wide constants
- Provides utility functions for common operations
- Contains shared error definitions
- Offers helper functions for time handling and formatting

## Main Files

### `constants.go`
Contains application constants and common errors.

**Cookie Names:**
```go
const (
    AccessTokenCookieName  = "access_token"
    RefreshTokenCookieName = "refresh_token"
)
```

**Common Errors:**
```go
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserNotFound       = errors.New("user not found")
    ... // Add more as needed
)
```

### `time_handling.go`
Contains time-related utility functions.

**Functions:**

**`GetUptimeString(startTime)`** - Formats application uptime
```go
// Returns human-readable uptime string
// Example: "2h 30m 45s"
func GetUptimeString(startTime time.Time) string
```

## How to Use

### Using Constants

```go
import "github.com/MonkyMars/PWS/lib"

// Set cookie with defined name
c.Cookie(&fiber.Cookie{
    Name:  lib.AccessTokenCookieName,
    Value: accessToken,
})

// Get cookie with defined name
token := c.Cookies(lib.AccessTokenCookieName)
```

### Using Error Constants

```go
import (
    "errors"
    "github.com/MonkyMars/PWS/lib"
)

// Check for specific error
if errors.Is(err, lib.ErrInvalidCredentials) {
    return response.Unauthorized(c, "Invalid email or password")
}

// Return standard error
if userNotFound {
    return lib.ErrUserNotFound
}
```

### Using Time Utilities

```go
import "github.com/MonkyMars/PWS/lib"

// Track application start time
var appStartTime = time.Now()

// Get formatted uptime
uptime := lib.GetUptimeString(appStartTime)
// Returns something like "1h 23m 45s"
```

## Constants Usage Examples

### Authentication Cookies
```go
// Setting authentication cookies
cookieService.SetCookie(c, lib.AccessTokenCookieName, accessToken, time.Hour)
cookieService.SetCookie(c, lib.RefreshTokenCookieName, refreshToken, 7*24*time.Hour)

// Reading authentication cookies
accessToken := c.Cookies(lib.AccessTokenCookieName)
refreshToken := c.Cookies(lib.RefreshTokenCookieName)

// Clearing authentication cookies
c.ClearCookie(lib.AccessTokenCookieName)
c.ClearCookie(lib.RefreshTokenCookieName)
```

## Error Handling Patterns

### Standard Error Checks
```go
func Login(email, password string) (*User, error) {
    user, err := getUserByEmail(email)
    if err != nil {
        return nil, lib.ErrUserNotFound
    }

    if !validatePassword(password, user.PasswordHash) {
        return nil, lib.ErrInvalidCredentials
    }

    return user, nil
}
```

### In HTTP Handlers
```go
func LoginHandler(c fiber.Ctx) error {
    user, err := authService.Login(email, password)
    if err != nil {
        if errors.Is(err, lib.ErrInvalidCredentials) {
            return response.Unauthorized(c, "Invalid email or password")
        }
        if errors.Is(err, lib.ErrUserNotFound) {
            return response.Unauthorized(c, "Invalid email or password")
        }
        return response.InternalServerError(c, "Login failed")
    }

    return response.Success(c, user)
}
```

## Time Utilities

### Application Uptime
```go
var appStartTime = time.Now()

func GetSystemHealth(c fiber.Ctx) error {
    return response.Success(c, types.HealthResponse{
        Status:            "ok",
        ApplicationUptime: lib.GetUptimeString(appStartTime),
        Services:          map[string]string{"api": "ok"},
    })
}
```

### Duration Formatting
The `GetUptimeString` function formats time durations in a human-readable way:
- Less than 1 minute: "45s"
- Less than 1 hour: "23m 45s"
- More than 1 hour: "2h 30m 45s"
- More than 1 day: "1d 5h 30m"

## Adding New Constants

When adding new constants, group them logically:

```go
// Cookie names
const (
    AccessTokenCookieName  = "access_token"
    RefreshTokenCookieName = "refresh_token"
    SessionCookieName      = "session_id"
)

// Request headers
const (
    AuthorizationHeader = "Authorization"
    ContentTypeHeader   = "Content-Type"
    UserAgentHeader     = "User-Agent"
)

// Default values
const (
    DefaultPageSize     = 10
    MaxPageSize        = 100
    DefaultCacheExpiry = 5 * time.Minute
)
```

## Adding New Errors

Follow the error naming convention:

```go
var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrUserNotFound       = errors.New("user not found")
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrTokenExpired       = errors.New("token expired")
    ErrInsufficientPermissions = errors.New("insufficient permissions")
)
```

## Best Practices

1. **Use descriptive names** - Constants should clearly indicate their purpose
2. **Group related constants** - Keep similar constants together
3. **Consistent naming** - Follow Go naming conventions (CamelCase for exported)
4. **Document purpose** - Add comments for complex constants
5. **Avoid magic strings** - Use constants instead of hardcoded strings
6. **Standard errors** - Define common errors once and reuse them
7. **Semantic grouping** - Organize constants and functions by their domain

## Security Considerations

1. **Cookie names** - Use consistent, secure cookie names
2. **Error messages** - Don't expose sensitive information in error constants
3. **Constants** - Don't include secrets or credentials in constants
4. **Error handling** - Use generic error messages for security-sensitive operations

## Testing

Constants and utilities can be tested:

```go
func TestErrorConstants(t *testing.T) {
    assert.Equal(t, "invalid credentials", lib.ErrInvalidCredentials.Error())
    assert.Equal(t, "user not found", lib.ErrUserNotFound.Error())
}

func TestCookieNames(t *testing.T) {
    assert.Equal(t, "access_token", lib.AccessTokenCookieName)
    assert.Equal(t, "refresh_token", lib.RefreshTokenCookieName)
}

func TestUptimeString(t *testing.T) {
    start := time.Now().Add(-2 * time.Hour)
    uptime := lib.GetUptimeString(start)
    assert.Contains(t, uptime, "h")
}
```

This package provides the foundational utilities and constants that make the codebase more maintainable and consistent.
