package lib

import (
	"errors"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// Common application errors
var (
	// Authentication & Authorization errors
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInvalidToken            = errors.New("invalid token")
	ErrExpiredToken            = errors.New("expired token")
	ErrTokenGeneration         = errors.New("error generating token")
	ErrGeneratingToken         = errors.New("error generating token") // Alias for backwards compatibility
	ErrTokenValidation         = errors.New("error validating token")
	ErrValidatingToken         = errors.New("error validating token") // Alias for backwards compatibility
	ErrTokenRefresh            = errors.New("failed to refresh token")
	ErrFailedToRefreshToken    = errors.New("failed to refresh token") // Alias for backwards compatibility
	ErrTokenDeletion           = errors.New("failed to delete token")
	ErrFailedToDeleteToken     = errors.New("failed to delete token") // Alias for backwards compatibility
	ErrUnauthorized            = errors.New("unauthorized access")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrTokenRevoked            = errors.New("token has been revoked")

	// User management errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrPasswordHashing   = errors.New("error hashing password")
	ErrHashingPassword   = errors.New("error hashing password") // Alias for backwards compatibility
	ErrUserCreation      = errors.New("error creating user")
	ErrCreateUser        = errors.New("error creating user") // Alias for backwards compatibility

	// Content management errors
	ErrFileNotFound   = errors.New("file not found")
	ErrFileUpload     = errors.New("file upload failed")
	ErrFileAccess     = errors.New("file access denied")
	ErrFolderNotFound = errors.New("folder not found")
	ErrFolderCreation = errors.New("folder creation failed")

	// Validation errors
	ErrInvalidInput     = errors.New("invalid input data")
	ErrMissingField     = errors.New("required field missing")
	ErrMissingParameter = errors.New("missing URL parameter")
	ErrInvalidFormat    = errors.New("invalid data format")
	ErrMissingFile      = errors.New("missing file in request")

	// Service errors
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrExternalService    = errors.New("external service error")
	ErrNotFound           = errors.New("resource not found")
)

// ErrorHandler provides centralized error handling with consistent responses
type ErrorHandler struct {
	logger *config.Logger
}

// NewErrorHandler creates a new error handler instance
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		logger: config.SetupLogger(),
	}
}

// HandleServiceError provides standardized error response patterns
// This function maps service layer errors to appropriate HTTP responses
func HandleServiceError(c fiber.Ctx, err error) error {
	handler := NewErrorHandler()
	return handler.Handle(c, err)
}

// Handle processes errors and returns appropriate HTTP responses
func (eh *ErrorHandler) Handle(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Log the error for debugging
	eh.logError(c, err)

	// Map specific errors to HTTP responses
	switch {
	// Authentication & Authorization errors (401)
	case errors.Is(err, ErrInvalidCredentials):
		return response.Unauthorized(c, "Invalid email or password")
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpiredToken):
		return response.Unauthorized(c, "Invalid or expired token")
	case errors.Is(err, ErrUnauthorized):
		return response.Unauthorized(c, "Authentication required")

	// Forbidden access (403)
	case errors.Is(err, ErrInsufficientPermissions):
		return response.Forbidden(c, "You do not have permission to perform this action")
	case errors.Is(err, ErrFileAccess):
		return response.Forbidden(c, "File access denied")

	// Not Found errors (404)
	case errors.Is(err, ErrUserNotFound):
		return response.NotFound(c, "User not found")
	case errors.Is(err, ErrFileNotFound):
		return response.NotFound(c, "File not found")
	case errors.Is(err, ErrFolderNotFound):
		return response.NotFound(c, "Folder not found")

	// Conflict errors (409)
	case errors.Is(err, ErrUserAlreadyExists):
		return response.Conflict(c, "User with this email already exists")
	case errors.Is(err, ErrUsernameTaken):
		return response.Conflict(c, "Username is already taken")

	// Bad Request errors (400)
	case errors.Is(err, ErrInvalidInput), errors.Is(err, ErrInvalidFormat), errors.Is(err, ErrMissingFile):
		return response.BadRequest(c, "Invalid input data")
	case errors.Is(err, ErrMissingField), errors.Is(err, ErrMissingParameter):
		return response.BadRequest(c, "Required field is missing")

	// Service Unavailable errors (503)
	case errors.Is(err, ErrServiceUnavailable):
		return response.ServiceUnavailable(c, "Service temporarily unavailable")

	// Token generation/management errors (500)
	case errors.Is(err, ErrTokenGeneration):
		return response.InternalServerError(c, "Failed to generate authentication token")
	case errors.Is(err, ErrTokenRefresh):
		return response.InternalServerError(c, "Failed to refresh token")

	// User management errors (500)
	case errors.Is(err, ErrPasswordHashing), errors.Is(err, ErrUserCreation):
		return response.InternalServerError(c, "User account creation failed")

	// File/Content management errors (500)
	case errors.Is(err, ErrFileUpload):
		return response.InternalServerError(c, "File upload failed")
	case errors.Is(err, ErrFolderCreation):
		return response.InternalServerError(c, "Folder creation failed")

	// Database/Infrastructure errors (500)
	case errors.Is(err, ErrDatabaseConnection):
		return response.InternalServerError(c, "Database connection error")
	case errors.Is(err, ErrExternalService):
		return response.InternalServerError(c, "External service error")

	// Default case for unknown errors (500)
	default:
		return response.InternalServerError(c, "An unexpected error occurred")
	}
}

// HandleValidationError handles request validation errors specifically
func HandleValidationError(c fiber.Ctx, err error, field string) error {
	handler := NewErrorHandler()
	handler.logError(c, err)

	return response.SendValidationError(c, []types.ValidationError{
		{
			Field:   field,
			Message: err.Error(),
		},
	})
}

// HandleAuthError handles authentication-specific errors with context
func HandleAuthError(c fiber.Ctx, err error, context string) error {
	handler := NewErrorHandler()
	handler.logger.AuditError("Authentication error", "context", context, "error", err.Error())

	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return response.Unauthorized(c, "Invalid email or password")
	case errors.Is(err, ErrTokenGeneration), errors.Is(err, ErrTokenRefresh):
		return response.InternalServerError(c, "Authentication service temporarily unavailable")
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpiredToken), errors.Is(err, ErrUnauthorized):
		return response.Unauthorized(c, "Invalid or expired token")
	default:
		return response.Unauthorized(c, "Authentication failed")
	}
}

// logError logs errors with request context for debugging
func (eh *ErrorHandler) logError(c fiber.Ctx, err error) {
	eh.logger.AuditError("Request error",
		"error", err.Error(),
		"method", c.Method(),
		"path", c.Path(),
		"ip", c.IP(),
	)
}
