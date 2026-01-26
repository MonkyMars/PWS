package lib

import (
	"errors"
	"log"

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
	ErrTokenReuse              = errors.New("possible token reuse detected")
	ErrInvalidClaims           = errors.New("invalid authentication claims")

	// User management errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrPasswordHashing   = errors.New("error hashing password")
	ErrHashingPassword   = errors.New("error hashing password") // Alias for backwards compatibility
	ErrUserCreation      = errors.New("error creating user")
	ErrCreateUser        = errors.New("error creating user") // Alias for backwards compatibility
	ErrPasswordMismatch  = errors.New("password and confirmation do not match")
	ErrWeakPassword      = errors.New("password does not meet strength requirements")

	// Content management errors
	ErrFileNotFound    = errors.New("file not found")
	ErrFileUpload      = errors.New("file upload failed")
	ErrFileAccess      = errors.New("file access denied")
	ErrFolderNotFound  = errors.New("folder not found")
	ErrFolderCreation  = errors.New("folder creation failed")
	ErrSubjectNotFound = errors.New("subject not found")
	ErrServiceNotFound = errors.New("service not found")

	// Validation errors
	ErrInvalidInput     = errors.New("invalid input data")
	ErrMissingField     = errors.New("required field missing")
	ErrInvalidFormat    = errors.New("invalid data format")
	ErrInvalidRequest   = errors.New("invalid request")
	ErrValidation       = errors.New("validation error")
	ErrMissingFile      = errors.New("Missing file(s)")
	ErrMissingParameter = errors.New("Missing file(s)")

	// Access control errors
	ErrForbidden = errors.New("forbidden access")

	// External service errors
	ErrNoLinkedAccount = errors.New("no linked account")

	// Service errors
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrExternalService    = errors.New("external service error")
	ErrWorkerUnavailable  = errors.New("worker unavailable")
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
// The message parameter is logged for developers but not exposed to users
func HandleServiceError(c fiber.Ctx, err error, message string) error {
	handler := NewErrorHandler()
	return handler.Handle(c, err, message)
}

// HandleServiceWarning logs a warning with context but continues with success response
// Use this for non-critical issues that should be monitored but don't block the operation
func HandleServiceWarning(c fiber.Ctx, message string, data ...any) {
	handler := NewErrorHandler()
	if handler.logger != nil {
		args := []any{"message", message, "method", c.Method(), "path", c.Path(), "ip", c.IP()}
		args = append(args, data...)
		handler.logger.Warn("Service warning", args...)
	}
}

// Handle processes errors and returns appropriate HTTP responses
func (eh *ErrorHandler) Handle(c fiber.Ctx, err error, message string) error {
	if err == nil {
		return nil
	}

	// Log the error with detailed message for developers
	eh.logErrorWithMessage(c, err, message)

	// Map specific errors to HTTP responses
	switch {
	// Authentication & Authorization errors (401)
	case errors.Is(err, ErrInvalidCredentials):
		return response.Unauthorized(c, "Invalid email or password")
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpiredToken):
		return response.Unauthorized(c, "Invalid or expired token")
	case errors.Is(err, ErrUnauthorized):
		return response.Unauthorized(c, "Authentication required")
	case errors.Is(err, ErrTokenReuse):
		return response.Unauthorized(c, "Invalid token")
	case errors.Is(err, ErrInvalidClaims):
		return response.Unauthorized(c, "Unauthorized")
	case errors.Is(err, ErrTokenRevoked):
		return response.Unauthorized(c, "Token has been revoked")

	// Forbidden access (403)
	case errors.Is(err, ErrInsufficientPermissions):
		return response.Forbidden(c, "You do not have permission to perform this action")
	case errors.Is(err, ErrFileAccess):
		return response.Forbidden(c, "File access denied")
	case errors.Is(err, ErrForbidden):
		return response.Forbidden(c, "Access denied")

	// Not Found errors (404)
	case errors.Is(err, ErrUserNotFound):
		return response.NotFound(c, "User not found")
	case errors.Is(err, ErrFileNotFound):
		return response.NotFound(c, "File not found")
	case errors.Is(err, ErrFolderNotFound):
		return response.NotFound(c, "Folder not found")
	case errors.Is(err, ErrSubjectNotFound):
		return response.NotFound(c, "Subject not found")
	case errors.Is(err, ErrServiceNotFound):
		return response.NotFound(c, "Service not found")
	case errors.Is(err, ErrNoLinkedAccount):
		return response.NotFound(c, "No linked account found")
	case errors.Is(err, ErrNotFound):
		return response.NotFound(c, "Resource not found")

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
	case errors.Is(err, ErrInvalidRequest):
		return response.BadRequest(c, "Invalid request")
	case errors.Is(err, ErrValidation):
		return response.BadRequest(c, "Validation failed")

	// Service Unavailable errors (503)
	case errors.Is(err, ErrServiceUnavailable):
		return response.ServiceUnavailable(c, "Service temporarily unavailable")
	case errors.Is(err, ErrWorkerUnavailable):
		return response.ServiceUnavailable(c, "Service temporarily unavailable")

	// Token generation/management errors (500)
	case errors.Is(err, ErrTokenGeneration):
		return response.InternalServerError(c, "Failed to generate authentication token")
	case errors.Is(err, ErrTokenRefresh):
		return response.InternalServerError(c, "Failed to refresh token")
	case errors.Is(err, ErrTokenDeletion):
		return response.InternalServerError(c, "Failed to revoke token")

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

// GetValidatedClaims extracts and validates authentication claims from context
// Returns the claims or an error that can be passed to HandleServiceError
func GetValidatedClaims(c fiber.Ctx) (*types.AuthClaims, error) {
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		return nil, ErrInvalidClaims
	}

	claims, ok := claimsInterface.(*types.AuthClaims)
	if !ok || claims == nil {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// logErrorWithMessage logs errors with detailed message and request context
func (eh *ErrorHandler) logErrorWithMessage(c fiber.Ctx, err error, message string) {
	if eh.logger != nil {
		eh.logger.AuditError(
			message,
			"error", err.Error(),
			"method", c.Method(),
			"path", c.Path(),
			"ip", c.IP(),
		)
	} else {
		log.Printf("Error: %s | %v, Method: %s, Path: %s", message, err, c.Method(), c.Path())
	}
}
