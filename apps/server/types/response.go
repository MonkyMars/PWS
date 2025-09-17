package types

import "time"

type Response struct {
	// Success indicates whether the request was processed successfully
	Success bool `json:"success"`
	// Message provides a human-readable description of the response
	Message string `json:"message"`
	// Data contains the response payload for successful requests (omitted if nil)
	Data any `json:"data,omitempty"`
	// Error contains detailed error information for failed requests (omitted if nil)
	Error *ErrorInfo `json:"error,omitempty"`
	// Meta contains additional metadata, particularly useful for pagination (omitted if nil)
	Meta *Meta `json:"meta,omitempty"`
	// Timestamp is the Unix timestamp in milliseconds when the response was generated
	Timestamp time.Time `json:"timestamp"`
}

// ErrorInfo contains detailed error information for failed requests.
// This structure provides structured error data that clients can use for
// error handling, user feedback, and debugging.
type ErrorInfo struct {
	// Code is a machine-readable error code for programmatic error handling
	Code string `json:"code"`
	// Message is a human-readable error description
	Message string `json:"message"`
	// Details contains additional contextual information about the error (omitted if empty)
	Details map[string]any `json:"details,omitempty"`
	// Field specifies which input field caused the error (useful for validation errors)
	Field string `json:"field,omitempty"`
}

// ValidationError represents field-specific validation errors.
// This structure is typically used within error details to provide
// specific information about validation failures.
type ValidationError struct {
	// Field is the name of the field that failed validation
	Field string `json:"field"`
	// Message describes why the validation failed
	Message string `json:"message"`
	// Value is the invalid value that was provided (omitted if empty)
	Value string `json:"value,omitempty"`
}

// Meta contains metadata for responses, especially useful for pagination.
// This structure provides clients with information needed to implement
// pagination controls and understand data set boundaries.
type Meta struct {
	// Page is the current page number (1-based)
	Page int `json:"page,omitempty"`
	// Limit is the maximum number of items per page
	Limit int `json:"limit,omitempty"`
	// Total is the total number of items across all pages
	Total int `json:"total,omitempty"`
	// TotalPages is the total number of pages available
	TotalPages int `json:"total_pages,omitempty"`
	// HasNext indicates if there are more pages after the current one
	HasNext bool `json:"has_next,omitempty"`
	// HasPrev indicates if there are pages before the current one
	HasPrev bool `json:"has_prev,omitempty"`
}

// PaginatedData wraps data with pagination metadata.
// This structure is used for list endpoints that support pagination.
type PaginatedData struct {
	// Items contains the array of data items for the current page
	Items []any `json:"items"`
	// Meta contains pagination metadata
	Meta *Meta `json:"meta"`
}
