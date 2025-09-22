package response

import (
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// Success sends a successful response with data using a default success message.
// This is a convenience function for common success responses.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - data: The data to include in the response
//
// Returns an error if the response cannot be sent.
func Success(c fiber.Ctx, data any) error {
	return NewResponse().
		Success("Request successful").
		WithData(data).
		Send(c, fiber.StatusOK)
}

// SuccessWithMessage sends a successful response with a custom message and data.
// This function allows for more specific success messages while including response data.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom success message
//   - data: The data to include in the response
//
// Returns an error if the response cannot be sent.
func SuccessWithMessage(c fiber.Ctx, message string, data any) error {
	return NewResponse().
		Success(message).
		WithData(data).
		Send(c, fiber.StatusOK)
}

// Created sends a 201 Created response for successful resource creation.
// This function should be used when a new resource has been successfully created.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - data: The created resource data
//
// Returns an error if the response cannot be sent.
func Created(c fiber.Ctx, data any) error {
	return NewResponse().
		Success("Resource created successfully").
		WithData(data).
		Send(c, fiber.StatusCreated)
}

// CreatedWithMessage sends a 201 Created response with a custom message.
// This function allows for specific creation messages while returning the created resource.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom creation success message
//   - data: The created resource data
//
// Returns an error if the response cannot be sent.
func CreatedWithMessage(c fiber.Ctx, message string, data any) error {
	return NewResponse().
		Success(message).
		WithData(data).
		Send(c, fiber.StatusCreated)
}

// Accepted sends a 202 Accepted response for requests that have been accepted for processing.
// This function should be used for asynchronous operations that have been queued or started.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Message describing the accepted operation
//
// Returns an error if the response cannot be sent.
func Accepted(c fiber.Ctx, message string) error {
	return NewResponse().
		Success(message).
		Send(c, fiber.StatusAccepted)
}

// NoContent sends a 204 No Content response for successful operations without response body.
// This function should be used for operations like DELETE that succeed but don't return data.
//
// Parameters:
//   - c: Fiber context for sending the response
//
// Returns an error if the response cannot be sent.
func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// OKWithMeta sends a successful response with both data and metadata.
// This function is useful for responses that need to include additional metadata.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - data: The response data
//   - meta: Additional metadata (often pagination information)
//
// Returns an error if the response cannot be sent.
func OKWithMeta(c fiber.Ctx, data any, meta *types.Meta) error {
	return NewResponse().
		Success("Request successful").
		WithData(data).
		WithMeta(meta).
		Send(c, fiber.StatusOK)
}

// Paginated sends a paginated response with automatically calculated pagination metadata.
// This function simplifies sending paginated responses by handling metadata calculation.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - items: Array of items for the current page
//   - page: Current page number (1-based)
//   - limit: Maximum number of items per page
//   - total: Total number of items across all pages
//
// Returns an error if the response cannot be sent.
func Paginated(c fiber.Ctx, items []any, page, limit, total int) error {
	meta := NewMeta(page, limit, total)
	paginatedData := NewPaginatedData(items, meta)

	return NewResponse().
		Success("Data retrieved successfully").
		WithData(paginatedData).
		WithMeta(meta).
		Send(c, fiber.StatusOK)
}

// PaginatedWithMessage sends a paginated response with a custom success message.
// This function combines pagination functionality with custom messaging.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Custom success message
//   - items: Array of items for the current page
//   - page: Current page number (1-based)
//   - limit: Maximum number of items per page
//   - total: Total number of items across all pages
//
// Returns an error if the response cannot be sent.
func PaginatedWithMessage(c fiber.Ctx, message string, items []any, page, limit, total int) error {
	meta := NewMeta(page, limit, total)
	paginatedData := NewPaginatedData(items, meta)

	return NewResponse().
		Success(message).
		WithData(paginatedData).
		WithMeta(meta).
		Send(c, fiber.StatusOK)
}

// Message sends a simple success message response without any data payload.
// This function is useful for operations that succeed but don't need to return data.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - message: Success message to send
//
// Returns an error if the response cannot be sent.
func Message(c fiber.Ctx, message string) error {
	return NewResponse().
		Success(message).
		Send(c, fiber.StatusOK)
}
