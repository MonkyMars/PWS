package response

import (
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)
// BadRequest sends a 400 Bad Request response for malformed or invalid requests.
// This function should be used when the client sends invalid request data.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Human-readable error message describing the bad request
//
// Returns an error if the response cannot be sent.
func BadRequest(c fiber.Ctx, message string) error {
	return NewResponse().
		Error(message).
		WithError(ErrCodeBadRequest, message).
		Send(c, fiber.StatusBadRequest)
}

// BadRequestWithDetails sends a 400 response with detailed error information.
// This function provides additional context about what made the request invalid.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Human-readable error message
//   - details: Map containing additional error context and details
//
// Returns an error if the response cannot be sent.
func BadRequestWithDetails(c fiber.Ctx, message string, details map[string]any) error {
	return NewResponse().
		Error(message).
		WithError(ErrCodeBadRequest, message).
		WithErrorDetails(details).
		Send(c, fiber.StatusBadRequest)
}

// Unauthorized sends a 401 Unauthorized response for authentication failures.
// This function should be used when the client needs to authenticate to access the resource.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom error message (uses default if empty)
//
// Returns an error if the response cannot be sent.
func Unauthorized(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Authentication required"
	}
	return NewResponse().
		Error(message).
		WithError(ErrCodeUnauthorized, message).
		Send(c, fiber.StatusUnauthorized)
}

// Forbidden sends a 403 Forbidden response for authorization failures.
// This function should be used when the authenticated user lacks permission for the resource.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom error message (uses default if empty)
//
// Returns an error if the response cannot be sent.
func Forbidden(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Access forbidden"
	}
	return NewResponse().
		Error(message).
		WithError(ErrCodeForbidden, message).
		Send(c, fiber.StatusForbidden)
}

// NotFound sends a 404 Not Found response when a requested resource doesn't exist.
// This function should be used when the requested resource cannot be found.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom error message (uses default if empty)
//
// Returns an error if the response cannot be sent.
func NotFound(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return NewResponse().
		Error(message).
		WithError(ErrCodeNotFound, message).
		Send(c, fiber.StatusNotFound)
}

// Conflict sends a 409 Conflict response for resource state conflicts.
// This function should be used when the request conflicts with the current resource state.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Error message describing the conflict
//
// Returns an error if the response cannot be sent.
func Conflict(c fiber.Ctx, message string) error {
	return NewResponse().
		Error(message).
		WithError(ErrCodeConflict, message).
		Send(c, fiber.StatusConflict)
}

// SendValidationError sends a 422 Unprocessable Entity response for validation errors.
// This function should be used when request data fails validation rules.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - errors: Array of ValidationError structs describing each validation failure
//
// Returns an error if the response cannot be sent.
func SendValidationError(c fiber.Ctx, errors []types.ValidationError) error {
	details := make(map[string]any)
	details["validation_errors"] = errors

	return NewResponse().
		Error("Validation failed").
		WithError(ErrCodeValidation, "The given data was invalid").
		WithErrorDetails(details).
		Send(c, fiber.StatusUnprocessableEntity)
}

// ValidationErrorSingle sends a 422 response for a single field validation error.
// This function provides a convenient way to report validation errors for individual fields.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - field: Name of the field that failed validation
//   - message: Validation error message
//   - value: The invalid value that was provided
//
// Returns an error if the response cannot be sent.
func ValidationErrorSingle(c fiber.Ctx, field, message, value string) error {
	return NewResponse().
		Error("Validation failed").
		WithError(ErrCodeValidation, message).
		WithField(field).
		WithErrorDetails(map[string]any{
			"field": field,
			"value": value,
		}).
		Send(c, fiber.StatusUnprocessableEntity)
}

// TooManyRequests sends a 429 Too Many Requests response for rate limiting.
// This function should be used when the client has exceeded rate limits.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom error message (uses default if empty)
//
// Returns an error if the response cannot be sent.
func TooManyRequests(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Too many requests"
	}
	return NewResponse().
		Error(message).
		WithError(ErrCodeTooManyReq, message).
		Send(c, fiber.StatusTooManyRequests)
}

// InternalServerError sends a 500 Internal Server Error response for server-side errors.
// This function should be used when an unexpected server error occurs.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom error message (uses default if empty)
//
// Returns an error if the response cannot be sent.
func InternalServerError(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return NewResponse().
		Error(message).
		WithError(ErrCodeInternal, message).
		Send(c, fiber.StatusInternalServerError)
}

// InternalServerErrorWithDetails sends a 500 response with detailed error information.
// This function provides additional context for debugging server errors.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom error message (uses default if empty)
//   - details: Map containing additional error context for debugging
//
// Returns an error if the response cannot be sent.
func InternalServerErrorWithDetails(c fiber.Ctx, message string, details map[string]any) error {
	if message == "" {
		message = "Internal server error"
	}
	return NewResponse().
		Error(message).
		WithError(ErrCodeInternal, message).
		WithErrorDetails(details).
		Send(c, fiber.StatusInternalServerError)
}

// ServiceUnavailable sends a 503 Service Unavailable response for temporary outages.
// This function should be used when the service is temporarily unable to handle requests.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom error message (uses default if empty)
//
// Returns an error if the response cannot be sent.
func ServiceUnavailable(c fiber.Ctx, message string) error {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	return NewResponse().
		Error(message).
		WithError(ErrCodeServiceUnavail, message).
		Send(c, fiber.StatusServiceUnavailable)
}

// CustomError sends a custom error response with the specified status code.
// This function provides flexibility for sending non-standard error responses.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - statusCode: HTTP status code to send
//   - code: Machine-readable error code
//   - message: Human-readable error message
//
// Returns an error if the response cannot be sent.
func CustomError(c fiber.Ctx, statusCode int, code, message string) error {
	return NewResponse().
		Error(message).
		WithError(code, message).
		Send(c, statusCode)
}

// CustomErrorWithDetails sends a custom error response with additional details.
// This function combines custom status codes with detailed error information.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - statusCode: HTTP status code to send
//   - code: Machine-readable error code
//   - message: Human-readable error message
//   - details: Map containing additional error context
//
// Returns an error if the response cannot be sent.
func CustomErrorWithDetails(c fiber.Ctx, statusCode int, code, message string, details map[string]any) error {
	return NewResponse().
		Error(message).
		WithError(code, message).
		WithErrorDetails(details).
		Send(c, statusCode)
}

// HandleError is a utility function to handle common Go errors with generic error responses.
// This function provides a convenient way to convert Go errors into HTTP error responses.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - err: The Go error to handle
//
// Returns an error if the response cannot be sent, or nil if no error was provided.
func HandleError(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	return InternalServerError(c, err.Error())
}