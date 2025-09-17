# Types Package

The `types` package defines shared data structures, interfaces, and type definitions used across the PWS server application. This package provides a centralized location for type definitions that ensure consistency and type safety throughout the codebase.

## Overview

This package contains:

- Data transfer objects (DTOs) for API requests and responses
- Database model structures
- Interface definitions for service contracts
- Shared enums and constants
- Validation tags and custom types

## Structure

Types are organized by their primary usage context:

```
types/
├── README.md          # This documentation
├── models.go          # Database model structures
├── requests.go        # API request DTOs
├── responses.go       # API response DTOs
├── interfaces.go      # Service and repository interfaces
├── enums.go          # Enumerated types and constants
└── validation.go     # Custom validation types and rules
```

## Data Models

### Database Models

Database models represent the structure of data stored in the database:

```go
package types

import (
    "time"
)

// User represents a user in the system.
type User struct {
    ID        string    `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Name      string    `json:"name" db:"name"`
    Role      UserRole  `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Note represents a note/document in the system.
type Note struct {
    ID          string    `json:"id" db:"id"`
    Title       string    `json:"title" db:"title"`
    Content     string    `json:"content" db:"content"`
    AuthorID    string    `json:"author_id" db:"author_id"`
    CategoryID  *string   `json:"category_id" db:"category_id"`
    IsPublic    bool      `json:"is_public" db:"is_public"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Category represents a note category.
type Category struct {
    ID          string    `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    Color       string    `json:"color" db:"color"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### Request DTOs

Request data transfer objects define the structure of incoming API requests:

```go
// CreateUserRequest defines the payload for creating a new user.
type CreateUserRequest struct {
    Email    string   `json:"email" validate:"required,email"`
    Name     string   `json:"name" validate:"required,min=2,max=100"`
    Password string   `json:"password" validate:"required,min=8"`
    Role     UserRole `json:"role" validate:"required"`
}

// UpdateUserRequest defines the payload for updating an existing user.
type UpdateUserRequest struct {
    Email *string   `json:"email,omitempty" validate:"omitempty,email"`
    Name  *string   `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
    Role  *UserRole `json:"role,omitempty"`
}

// CreateNoteRequest defines the payload for creating a new note.
type CreateNoteRequest struct {
    Title      string  `json:"title" validate:"required,min=1,max=200"`
    Content    string  `json:"content" validate:"required"`
    CategoryID *string `json:"category_id,omitempty" validate:"omitempty,uuid"`
    IsPublic   *bool   `json:"is_public,omitempty"`
}

// UpdateNoteRequest defines the payload for updating an existing note.
type UpdateNoteRequest struct {
    Title      *string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
    Content    *string `json:"content,omitempty"`
    CategoryID *string `json:"category_id,omitempty" validate:"omitempty,uuid"`
    IsPublic   *bool   `json:"is_public,omitempty"`
}

// LoginRequest defines the payload for user authentication.
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```

### Response DTOs

Response data transfer objects define the structure of API responses:

```go
// UserResponse represents a user in API responses.
type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Role      UserRole  `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

// NoteResponse represents a note in API responses.
type NoteResponse struct {
    ID         string           `json:"id"`
    Title      string           `json:"title"`
    Content    string           `json:"content"`
    Author     UserResponse     `json:"author"`
    Category   *CategoryResponse `json:"category,omitempty"`
    IsPublic   bool             `json:"is_public"`
    CreatedAt  time.Time        `json:"created_at"`
    UpdatedAt  time.Time        `json:"updated_at"`
}

// CategoryResponse represents a category in API responses.
type CategoryResponse struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Color       string    `json:"color"`
    CreatedAt   time.Time `json:"created_at"`
}

// AuthResponse represents authentication response data.
type AuthResponse struct {
    User         UserResponse `json:"user"`
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    ExpiresIn    int64        `json:"expires_in"`
}
```

## Interfaces

Service and repository interfaces define contracts for business logic and data access:

```go
// UserService defines the interface for user-related business logic.
type UserService interface {
    Create(req CreateUserRequest) (*User, error)
    GetByID(id string) (*User, error)
    GetByEmail(email string) (*User, error)
    Update(id string, req UpdateUserRequest) (*User, error)
    Delete(id string) error
    List(offset, limit int) ([]*User, int, error)
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
    Create(user *User) error
    GetByID(id string) (*User, error)
    GetByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id string) error
    List(offset, limit int) ([]*User, int, error)
}

// NoteService defines the interface for note-related business logic.
type NoteService interface {
    Create(authorID string, req CreateNoteRequest) (*Note, error)
    GetByID(id string) (*Note, error)
    Update(id, authorID string, req UpdateNoteRequest) (*Note, error)
    Delete(id, authorID string) error
    ListByAuthor(authorID string, offset, limit int) ([]*Note, int, error)
    ListPublic(offset, limit int) ([]*Note, int, error)
}

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
    Login(req LoginRequest) (*AuthResponse, error)
    RefreshToken(refreshToken string) (*AuthResponse, error)
    ValidateToken(token string) (*User, error)
    Logout(userID string) error
}
```

## Enumerations

Enumerated types provide type-safe constants:

```go
// UserRole defines the possible user roles in the system.
type UserRole string

const (
    UserRoleAdmin     UserRole = "admin"
    UserRoleModerator UserRole = "moderator"
    UserRoleUser      UserRole = "user"
)

// String returns the string representation of the UserRole.
func (r UserRole) String() string {
    return string(r)
}

// IsValid checks if the UserRole is a valid value.
func (r UserRole) IsValid() bool {
    switch r {
    case UserRoleAdmin, UserRoleModerator, UserRoleUser:
        return true
    default:
        return false
    }
}

// NoteStatus defines the possible states of a note.
type NoteStatus string

const (
    NoteStatusDraft     NoteStatus = "draft"
    NoteStatusPublished NoteStatus = "published"
    NoteStatusArchived  NoteStatus = "archived"
)

// SortOrder defines the sorting direction for queries.
type SortOrder string

const (
    SortOrderAsc  SortOrder = "asc"
    SortOrderDesc SortOrder = "desc"
)
```

## Custom Types

Custom types provide additional type safety and validation:

```go
// Email represents a validated email address.
type Email string

// NewEmail creates a new Email from a string after validation.
func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return "", errors.New("invalid email format")
    }
    return Email(email), nil
}

// String returns the string representation of the Email.
func (e Email) String() string {
    return string(e)
}

// UUID represents a validated UUID string.
type UUID string

// NewUUID creates a new UUID from a string after validation.
func NewUUID(id string) (UUID, error) {
    if !isValidUUID(id) {
        return "", errors.New("invalid UUID format")
    }
    return UUID(id), nil
}

// String returns the string representation of the UUID.
func (u UUID) String() string {
    return string(u)
}
```

## Validation Tags

The package uses struct tags for input validation:

```go
// Common validation tags used:
// - required: Field is required
// - omitempty: Skip validation if field is empty
// - min/max: String length or numeric range validation
// - email: Email format validation
// - uuid: UUID format validation
// - oneof: Value must be one of specified options

type ExampleRequest struct {
    Name     string   `json:"name" validate:"required,min=2,max=100"`
    Email    string   `json:"email" validate:"required,email"`
    Age      int      `json:"age" validate:"min=18,max=120"`
    Role     UserRole `json:"role" validate:"required,oneof=admin moderator user"`
    Optional *string  `json:"optional,omitempty" validate:"omitempty,min=1"`
}
```

## Error Types

Package-specific error types for consistent error handling:

```go
// ServiceError represents errors from the service layer.
type ServiceError struct {
    Code    string
    Message string
    Cause   error
}

func (e ServiceError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

// Common service errors
var (
    ErrUserNotFound     = ServiceError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrUserExists       = ServiceError{Code: "USER_EXISTS", Message: "User already exists"}
    ErrInvalidPassword  = ServiceError{Code: "INVALID_PASSWORD", Message: "Invalid password"}
    ErrNoteNotFound     = ServiceError{Code: "NOTE_NOT_FOUND", Message: "Note not found"}
    ErrUnauthorized     = ServiceError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
)
```

## Usage Examples

### Creating and Using Types

```go
package main

import (
    "github.com/MonkyMars/PWS/types"
)

func exampleUsage() {
    // Create a user request
    req := types.CreateUserRequest{
        Email:    "user@example.com",
        Name:     "John Doe",
        Password: "securepassword",
        Role:     types.UserRoleUser,
    }

    // Validate the request (using validation library)
    if err := validator.Struct(req); err != nil {
        // Handle validation error
    }

    // Create user model
    user := &types.User{
        ID:        generateUUID(),
        Email:     req.Email,
        Name:      req.Name,
        Role:      req.Role,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Convert to response DTO
    response := types.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Name:      user.Name,
        Role:      user.Role,
        CreatedAt: user.CreatedAt,
    }
}
```

## Best Practices

1. **Separation of Concerns**: Keep request, response, and model types separate
2. **Validation**: Use struct tags for input validation
3. **Immutability**: Use pointer fields for optional update fields
4. **Type Safety**: Use custom types for domain-specific values
5. **Documentation**: Document all public types and their fields
6. **Consistency**: Follow consistent naming and structure patterns
7. **Interfaces**: Define clear interfaces for service contracts

## Dependencies

- `time` - Standard library for time handling
- Validation library (e.g., `github.com/go-playground/validator/v10`)
- Database drivers as needed for struct tags
