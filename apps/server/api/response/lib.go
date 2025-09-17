package response

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MonkyMars/PWS/types"

	"github.com/gofiber/fiber/v3"
)

// ParsePaginationParams extracts pagination parameters from query string with default values.
// This function provides a standardized way to handle pagination across all endpoints.
// It enforces reasonable defaults and limits to prevent abuse.
//
// Default values: page=1, limit=10, maximum limit=100
//
// Parameters:
//   - c: Fiber context containing the query parameters
//
// Returns:
//   - page: The requested page number (minimum 1)
//   - limit: The requested page size (minimum 1, maximum 100)
//   - err: Error if the parameters are invalid
func ParsePaginationParams(c fiber.Ctx) (page, limit int, err error) {
	// Default values
	page = 1
	limit = 10

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else if err != nil {
			return 0, 0, fmt.Errorf("invalid page parameter: %s", pageStr)
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		} else if err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter: %s", limitStr)
		} else if l > 100 {
			limit = 100 // Cap at 100 items per page
		}
	}

	return page, limit, nil
}

// ParsePaginationParamsWithDefaults extracts pagination parameters with custom defaults and limits.
// This function allows for flexible pagination configuration per endpoint.
//
// Parameters:
//   - c: Fiber context containing the query parameters
//   - defaultPage: Default page number if not specified
//   - defaultLimit: Default page size if not specified
//   - maxLimit: Maximum allowed page size
//
// Returns:
//   - page: The requested page number (minimum 1)
//   - limit: The requested page size (minimum 1, maximum maxLimit)
//   - err: Error if the parameters are invalid
func ParsePaginationParamsWithDefaults(c fiber.Ctx, defaultPage, defaultLimit, maxLimit int) (page, limit int, err error) {
	page = defaultPage
	limit = defaultLimit

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else if err != nil {
			return 0, 0, fmt.Errorf("invalid page parameter: %s", pageStr)
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= maxLimit {
			limit = l
		} else if err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter: %s", limitStr)
		} else if l > maxLimit {
			limit = maxLimit
		}
	}

	return page, limit, nil
}

// CalculateOffset calculates the database offset value from page number and limit.
// This function provides a standard way to convert pagination parameters to database offset.
//
// Parameters:
//   - page: Page number (1-based)
//   - limit: Number of items per page
//
// Returns the offset value for database queries (0-based).
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// ValidateJSON validates and parses JSON request body into the provided struct.
// This function provides a standardized way to handle JSON parsing with error handling.
//
// Parameters:
//   - c: Fiber context containing the request body
//   - v: Pointer to the structure to parse the JSON into
//
// Returns an error if the JSON is invalid or cannot be parsed.
func ValidateJSON(c fiber.Ctx, v any) error {
	if err := c.Bind().JSON(v); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	return nil
}

// GetUserID extracts the authenticated user ID from the request context.
// This function assumes that authentication middleware has set the user ID in the context.
//
// Parameters:
//   - c: Fiber context containing user information
//
// Returns:
//   - string: The user ID
//   - error: Error if user ID is not found or invalid
func GetUserID(c fiber.Ctx) (string, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return "", fmt.Errorf("user ID not found in context")
	}

	if id, ok := userID.(string); ok {
		return id, nil
	}

	return "", fmt.Errorf("invalid user ID format")
}

// GetUserRole extracts the authenticated user's role from the request context.
// This function assumes that authentication middleware has set the user role in the context.
//
// Parameters:
//   - c: Fiber context containing user information
//
// Returns:
//   - string: The user role
//   - error: Error if user role is not found or invalid
func GetUserRole(c fiber.Ctx) (string, error) {
	userRole := c.Locals("user_role")
	if userRole == nil {
		return "", fmt.Errorf("user role not found in context")
	}

	if role, ok := userRole.(string); ok {
		return role, nil
	}

	return "", fmt.Errorf("invalid user role format")
}

// PrettyJSON returns a pretty-printed JSON string for debugging purposes.
// This function is useful for logging or debugging API responses and requests.
//
// Parameters:
//   - v: Any value to be serialized to pretty JSON
//
// Returns a formatted JSON string or error message if serialization fails.
func PrettyJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling JSON: %v", err)
	}
	return string(b)
}

// BuildValidationErrors creates a slice of ValidationError from a map of field errors.
// This function provides a convenient way to convert simple error maps to structured validation errors.
//
// Parameters:
//   - errors: Map where keys are field names and values are error messages
//
// Returns a slice of ValidationError structs.
func BuildValidationErrors(errors map[string]string) []types.ValidationError {
	var validationErrors []types.ValidationError
	for field, message := range errors {
		validationErrors = append(validationErrors, types.ValidationError{
			Field:   field,
			Message: message,
		})
	}
	return validationErrors
}

// SendValidationErrors is a helper function to quickly send validation errors from a map.
// This function combines BuildValidationErrors and SendValidationError for convenience.
//
// Parameters:
//   - c: Fiber context for sending the response
//   - errors: Map of field validation errors
//
// Returns an error if the response cannot be sent.
func SendValidationErrors(c fiber.Ctx, errors map[string]string) error {
	validationErrors := BuildValidationErrors(errors)
	return SendValidationError(c, validationErrors)
}

// QuickPaginate is a helper function for common pagination scenarios.
// This function handles pagination parameter parsing and response formatting in one call.
//
// Parameters:
//   - c: Fiber context containing pagination parameters
//   - items: Array of items for the current page
//   - totalCount: Total number of items across all pages
//
// Returns an error if pagination parameters are invalid or response cannot be sent.
func QuickPaginate(c fiber.Ctx, items []any, totalCount int) error {
	page, limit, err := ParsePaginationParams(c)
	if err != nil {
		return BadRequest(c, err.Error())
	}

	return Paginated(c, items, page, limit, totalCount)
}

// ConvertToInterfaceSlice converts typed slices to []any for use with pagination functions.
// This function handles common slice types and provides a fallback for other types.
//
// Parameters:
//   - slice: Any slice type to convert
//
// Returns a slice of any interfaces containing the original data.
func ConvertToInterfaceSlice(slice any) []any {
	switch s := slice.(type) {
	case []any:
		return s
	case []string:
		result := make([]any, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result
	case []int:
		result := make([]any, len(s))
		for i, v := range s {
			result[i] = v
		}
		return result
	default:
		// Use reflection for other types if needed
		return []any{slice}
	}
}

// ErrorWithContext adds request context information to error responses.
// This function enriches error responses with debugging information like request ID and path.
//
// Parameters:
//   - c: Fiber context containing request information
//   - statusCode: HTTP status code to send
//   - message: Error message
//   - context: Additional context information
//
// Returns an error if the response cannot be sent.
func ErrorWithContext(c fiber.Ctx, statusCode int, message string, context map[string]any) error {
	// Add request context
	context["request_id"] = c.Get("X-Request-ID")
	context["path"] = c.Path()
	context["method"] = c.Method()

	return CustomErrorWithDetails(c, statusCode, ErrCodeInternal, message, context)
}

// SuccessWithContext adds request context information to success responses.
// This function can be used to enrich success responses with additional metadata.
//
// Parameters:
//   - c: Fiber context containing request information
//   - data: Response data
//   - context: Additional context information
//
// Returns an error if the response cannot be sent.
func SuccessWithContext(c fiber.Ctx, data any, context map[string]any) error {
	// Add request context to meta
	meta := &types.Meta{}
	if context != nil {
		// Convert context to meta format if needed
		// This is a simple implementation, you might want to structure it differently
	}

	return NewResponse().
		Success("Request successful").
		WithData(data).
		WithMeta(meta).
		Send(c, fiber.StatusOK)
}