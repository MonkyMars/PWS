// Package response provides standardized HTTP response structures and utilities for the PWS API.
// This package implements a consistent response format across all API endpoints, including
// success responses, error handling, pagination support, and response building utilities.
//
// The package follows REST API best practices by providing:
//   - Consistent response structure across all endpoints
//   - Detailed error information with error codes and context
//   - Pagination metadata for list endpoints
//   - Fluent response builder pattern for easy response construction
//   - Standardized error codes for common scenarios
//
// Example usage:
//
//	// Success response with data
//	return response.NewResponse().
//		Success("User created successfully").
//		WithData(user).
//		Send(c, fiber.StatusCreated)
//
//	// Error response with details
//	return response.NewResponse().
//		Error("Validation failed").
//		WithError(response.ErrCodeValidation, "Invalid input data").
//		WithField("email").
//		Send(c, fiber.StatusBadRequest)
package response

import (
	"time"

	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// Response represents the standard API response structure used across all endpoints.
// This structure ensures consistency in API responses and provides clients with
// predictable response formats for both success and error scenarios.

// ResponseBuilder helps build responses fluently using the builder pattern.
// This provides a convenient and readable way to construct API responses
// with method chaining.
type ResponseBuilder struct {
	response *types.Response
}

// NewResponse creates a new response builder with initialized timestamp.
// This is the entry point for building API responses using the fluent builder pattern.
//
// Returns a new ResponseBuilder instance ready for method chaining.
func NewResponse() *ResponseBuilder {
	return &ResponseBuilder{
		response: &types.Response{
			Timestamp: time.Now(),
		},
	}
}

// Success sets the response as successful with the provided message.
// This method should be used for successful operations.
//
// Parameters:
//   - message: Human-readable success message
//
// Returns the ResponseBuilder for method chaining.
func (rb *ResponseBuilder) Success(message string) *ResponseBuilder {
	rb.response.Success = true
	rb.response.Message = message
	return rb
}

// Error sets the response as an error with the provided message.
// This method should be used for failed operations.
//
// Parameters:
//   - message: Human-readable error message
//
// Returns the ResponseBuilder for method chaining.
func (rb *ResponseBuilder) Error(message string) *ResponseBuilder {
	rb.response.Success = false
	rb.response.Message = message
	return rb
}

// WithData adds data payload to the response.
// This method is typically used with successful responses to include the response data.
//
// Parameters:
//   - data: The data to include in the response (can be any serializable type)
//
// Returns the ResponseBuilder for method chaining.
func (rb *ResponseBuilder) WithData(data any) *ResponseBuilder {
	rb.response.Data = data
	return rb
}

// WithError adds detailed error information to the response.
// This method provides structured error details for better error handling.
//
// Parameters:
//   - code: Machine-readable error code for programmatic handling
//   - message: Human-readable error description
//
// Returns the ResponseBuilder for method chaining.
func (rb *ResponseBuilder) WithError(code, message string) *ResponseBuilder {
	rb.response.Error = &types.ErrorInfo{
		Code:    code,
		Message: message,
	}
	return rb
}

// WithErrorDetails adds additional error context and details.
// This method allows adding structured error details for debugging and context.
//
// Parameters:
//   - details: Map containing additional error context
//
// Returns the ResponseBuilder for method chaining.
func (rb *ResponseBuilder) WithErrorDetails(details map[string]any) *ResponseBuilder {
	if rb.response.Error == nil {
		rb.response.Error = &types.ErrorInfo{}
	}
	rb.response.Error.Details = details
	return rb
}

// WithField adds field information to error responses.
// This method is particularly useful for validation errors to indicate which field caused the error.
//
// Parameters:
//   - field: The name of the field that caused the error
//
// Returns the ResponseBuilder for method chaining.
func (rb *ResponseBuilder) WithField(field string) *ResponseBuilder {
	if rb.response.Error == nil {
		rb.response.Error = &types.ErrorInfo{}
	}
	rb.response.Error.Field = field
	return rb
}

// WithMeta adds metadata to the response.
// This method is typically used for pagination information in list responses.
//
// Parameters:
//   - meta: Metadata structure containing pagination or other response metadata
//
// Returns the ResponseBuilder for method chaining.
func (rb *ResponseBuilder) WithMeta(meta *types.Meta) *ResponseBuilder {
	rb.response.Meta = meta
	return rb
}

// Build returns the final constructed response.
// This method completes the building process and returns the Response struct.
//
// Returns the constructed Response instance.
func (rb *ResponseBuilder) Build() *types.Response {
	return rb.response
}

// Send sends the response with the specified HTTP status code.
// This method completes the response building and sends it to the client.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - statusCode: HTTP status code to send with the response
//
// Returns an error if the response cannot be sent.
func (rb *ResponseBuilder) Send(c fiber.Ctx, statusCode int) error {
	return c.Status(statusCode).JSON(rb.response)
}

// Common error codes used throughout the application.
// These constants provide standardized error codes for consistent error handling.
const (
	// ErrCodeValidation indicates request validation errors
	ErrCodeValidation = "VALIDATION_ERROR"
	// ErrCodeNotFound indicates that the requested resource was not found
	ErrCodeNotFound = "NOT_FOUND"
	// ErrCodeUnauthorized indicates authentication is required
	ErrCodeUnauthorized = "UNAUTHORIZED"
	// ErrCodeForbidden indicates the user lacks permission for the requested action
	ErrCodeForbidden = "FORBIDDEN"
	// ErrCodeConflict indicates a conflict with the current resource state
	ErrCodeConflict = "CONFLICT"
	// ErrCodeInternal indicates an internal server error
	ErrCodeInternal = "INTERNAL_ERROR"
	// ErrCodeBadRequest indicates malformed or invalid request data
	ErrCodeBadRequest = "BAD_REQUEST"
	// ErrCodeTooManyReq indicates rate limiting is in effect
	ErrCodeTooManyReq = "TOO_MANY_REQUESTS"
	// ErrCodeServiceUnavail indicates the service is temporarily unavailable
	ErrCodeServiceUnavail = "SERVICE_UNAVAILABLE"
)

// NewMeta creates pagination metadata based on current page, limit, and total count.
// This function calculates pagination values including total pages and navigation flags.
//
// Parameters:
//   - page: Current page number (1-based)
//   - limit: Maximum number of items per page
//   - total: Total number of items across all pages
//
// Returns a Meta struct with calculated pagination information.
func NewMeta(page, limit, total int) *types.Meta {
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	return &types.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// NewPaginatedData creates a paginated data structure combining items with metadata.
// This function provides a convenient way to package paginated results with their metadata.
//
// Parameters:
//   - items: Array of data items for the current page
//   - meta: Pagination metadata
//
// Returns a PaginatedData struct containing both items and pagination metadata.
func NewPaginatedData(items []any, meta *types.Meta) *types.PaginatedData {
	return &types.PaginatedData{
		Items: items,
		Meta:  meta,
	}
}
