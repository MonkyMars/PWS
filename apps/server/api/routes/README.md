# Routes Package

The `routes` package contains HTTP route handlers and endpoint definitions for the PWS API. This package organizes route handlers by feature or resource type, providing a clean separation of concerns for different API endpoints.

## Overview

This package follows RESTful API design principles and implements route handlers that:

- Handle HTTP requests and responses
- Validate request data
- Interact with service layers for business logic
- Return standardized API responses
- Implement proper error handling

## Structure (Example)

Route handlers are organized by resource or feature:

```
routes/
├── README.md          # This documentation
├── health.go          # Health check endpoints (if implemented)
├── users.go           # User management endpoints (if implemented)
├── auth.go            # Authentication endpoints (if implemented)
└── api.go             # General API routes (if implemented)
```

## Route Handler Pattern

All route handlers in this package follow a consistent pattern:

```go
package routes

import (
    "github.com/MonkyMars/PWS/api/response"
    "github.com/gofiber/fiber/v3"
)

// ResourceHandler handles HTTP requests for a specific resource.
// It validates input, processes the request through service layers,
// and returns a standardized API response.
//
// Parameters:
//   - c: Fiber context containing request and response objects
//
// Returns an error if the response cannot be sent.
func ResourceHandler(c fiber.Ctx) error {
    // 1. Extract and validate request parameters
    id := c.Params("id")
    if id == "" {
        return response.BadRequest(c, "Resource ID is required")
    }

    // 2. Parse and validate request body (if needed)
    var req ResourceRequest
    if err := response.ValidateJSON(c, &req); err != nil {
        return response.BadRequest(c, "Invalid request format")
    }

    // 3. Call service layer for business logic
    result, err := resourceService.Process(id, req)
    if err != nil {
        return response.InternalServerError(c, "Failed to process request")
    }

    // 4. Return standardized response
    return response.Success(c, result)
}
```

## Route Registration

Routes are registered through setup functions that are called from the main router:

```go
package routes

import (
    "github.com/gofiber/fiber/v3"
)

// SetupResourceRoutes registers all routes for a specific resource.
// This function should be called from the main router setup.
//
// Parameters:
//   - app: Fiber application instance to register routes on
func SetupResourceRoutes(app *fiber.App) {
    // Create route group
    api := app.Group("/api/v1")
    resource := api.Group("/resource")

    // Register CRUD endpoints
    resource.Get("/", ListResources)
    resource.Get("/:id", GetResource)
    resource.Post("/", CreateResource)
    resource.Put("/:id", UpdateResource)
    resource.Delete("/:id", DeleteResource)
}
```

## Request/Response Types

Each route handler should define clear request and response types:

```go
// Request types for input validation
type CreateResourceRequest struct {
    Name        string `json:"name" validate:"required,min=1,max=100"`
    Description string `json:"description" validate:"max=500"`
    Category    string `json:"category" validate:"required"`
}

type UpdateResourceRequest struct {
    Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
    Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
    Category    *string `json:"category,omitempty"`
}

// Response types for output formatting
type ResourceResponse struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

## Example Implementation

### Basic CRUD Operations

```go
package routes

import (
    "github.com/MonkyMars/PWS/api/response"
    "github.com/MonkyMars/PWS/services"
    "github.com/gofiber/fiber/v3"
)

// ListResources handles GET /api/v1/resources
func ListResources(c fiber.Ctx) error {
    // Parse pagination parameters
    page, limit, err := response.ParsePaginationParams(c)
    if err != nil {
        return response.BadRequest(c, err.Error())
    }

    // Get resources from service
    resources, total, err := services.ResourceService.List(page, limit)
    if err != nil {
        return response.InternalServerError(c, "Failed to fetch resources")
    }

    // Convert to interface slice for pagination
    items := response.ConvertToInterfaceSlice(resources)

    // Return paginated response
    return response.Paginated(c, items, page, limit, total)
}

// GetResource handles GET /api/v1/resources/:id
func GetResource(c fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return response.BadRequest(c, "Resource ID is required")
    }

    resource, err := services.ResourceService.GetByID(id)
    if err != nil {
        if errors.Is(err, services.ErrResourceNotFound) {
            return response.NotFound(c, "Resource not found")
        }
        return response.InternalServerError(c, "Failed to fetch resource")
    }

    return response.Success(c, resource)
}

// CreateResource handles POST /api/v1/resources
func CreateResource(c fiber.Ctx) error {
    var req CreateResourceRequest
    if err := response.ValidateJSON(c, &req); err != nil {
        return response.BadRequest(c, "Invalid request format")
    }

    // Validate request
    if errors := validateCreateResource(req); len(errors) > 0 {
        return response.SendValidationErrors(c, errors)
    }

    // Create resource
    resource, err := services.ResourceService.Create(req)
    if err != nil {
        return response.InternalServerError(c, "Failed to create resource")
    }

    return response.Created(c, resource)
}

// UpdateResource handles PUT /api/v1/resources/:id
func UpdateResource(c fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return response.BadRequest(c, "Resource ID is required")
    }

    var req UpdateResourceRequest
    if err := response.ValidateJSON(c, &req); err != nil {
        return response.BadRequest(c, "Invalid request format")
    }

    resource, err := services.ResourceService.Update(id, req)
    if err != nil {
        if errors.Is(err, services.ErrResourceNotFound) {
            return response.NotFound(c, "Resource not found")
        }
        return response.InternalServerError(c, "Failed to update resource")
    }

    return response.Success(c, resource)
}

// DeleteResource handles DELETE /api/v1/resources/:id
func DeleteResource(c fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return response.BadRequest(c, "Resource ID is required")
    }

    err := services.ResourceService.Delete(id)
    if err != nil {
        if errors.Is(err, services.ErrResourceNotFound) {
            return response.NotFound(c, "Resource not found")
        }
        return response.InternalServerError(c, "Failed to delete resource")
    }

    return response.NoContent(c)
}
```

### Authentication-Protected Routes

```go
// ProtectedResource handles authenticated requests
func ProtectedResource(c fiber.Ctx) error {
    // Extract user from context (set by auth middleware)
    userID, err := response.GetUserID(c)
    if err != nil {
        return response.Unauthorized(c, "Authentication required")
    }

    role, err := response.GetUserRole(c)
    if err != nil {
        return response.Forbidden(c, "Invalid user role")
    }

    // Check permissions
    if !hasPermission(role, "read:resources") {
        return response.Forbidden(c, "Insufficient permissions")
    }

    // Process request...
    return response.Success(c, data)
}
```

## Validation

Request validation should be handled consistently:

```go
func validateCreateResource(req CreateResourceRequest) map[string]string {
    errors := make(map[string]string)

    if req.Name == "" {
        errors["name"] = "Name is required"
    } else if len(req.Name) > 100 {
        errors["name"] = "Name must be less than 100 characters"
    }

    if req.Category == "" {
        errors["category"] = "Category is required"
    }

    return errors
}
```

## Error Handling

Use the response package for consistent error handling:

```go
// Bad request for client errors
return response.BadRequest(c, "Invalid input")

// Not found for missing resources
return response.NotFound(c, "Resource not found")

// Internal server error for unexpected errors
return response.InternalServerError(c, "Something went wrong")

// Validation errors for form validation
return response.SendValidationErrors(c, validationErrors)
```

## Best Practices

1. **Validation**: Always validate input parameters and request bodies
2. **Error Handling**: Use appropriate HTTP status codes and error messages
3. **Pagination**: Implement pagination for list endpoints
4. **Authentication**: Check authentication and authorization where required
5. **Logging**: Log important events and errors for debugging
6. **Testing**: Write unit tests for route handlers
7. **Documentation**: Document API endpoints with clear examples

## Testing

```go
package routes_test

import (
    "testing"
    "github.com/MonkyMars/PWS/api/routes"
    "github.com/gofiber/fiber/v3"
)

func TestGetResource(t *testing.T) {
    app := fiber.New()
    routes.SetupResourceRoutes(app)

    // Test successful get
    // Test not found case
    // Test invalid ID format
}
```

## Dependencies

- `github.com/gofiber/fiber/v3` - Web framework for HTTP handling
- `github.com/MonkyMars/PWS/api/response` - Standardized response utilities
- `github.com/MonkyMars/PWS/services` - Business logic layer
- Validation libraries as needed

The routes package contains HTTP route handlers for the PWS server. Route handlers are functions that process incoming HTTP requests and generate appropriate responses using the response package.

## Structure

- Add a new file for each resource or feature of the system

## Handler Pattern

All route handlers follow a consistent pattern:

1. Accept a Fiber context parameter
2. Extract and validate request data
3. Perform business logic operations
4. Return standardized responses using the response package

```go
func HandlerName(c fiber.Ctx) error {
    // Extract request parameters
    id := c.Params("id")

    // Validate input
    if id == "" {
        return response.BadRequest(c, "ID is required")
    }

    // Perform business logic
    data, err := businessLogic(id)
    if err != nil {
        return response.InternalServerError(c, err.Error())
    }

    // Return success response
    return response.Success(c, data)
}
```

## Common Handler Patterns

### GET Handlers

- Extract query parameters for filtering and pagination
- Retrieve data from storage layer
- Return data with appropriate metadata

### POST Handlers

- Validate JSON request body
- Create new resources
- Return created resource with 201 status

### PUT/PATCH Handlers

- Validate request body and parameters
- Update existing resources
- Return updated resource or 204 No Content

### DELETE Handlers

- Validate resource exists
- Perform deletion
- Return 204 No Content or confirmation message

## Request Validation

Use the response package utilities for common validation tasks:

```go
// Validate JSON body
var req CreateNoteRequest
if err := response.ValidateJSON(c, &req); err != nil {
    return response.BadRequest(c, "Invalid JSON")
}

// Validate pagination parameters
page, limit, err := response.ParsePaginationParams(c)
if err != nil {
    return response.BadRequest(c, err.Error())
}
```

## Error Handling

Use appropriate response functions for different error scenarios:

```go
// Resource not found
if note == nil {
    return response.NotFound(c, "Note not found")
}

// Validation errors
if errors := validateNote(req); len(errors) > 0 {
    return response.SendValidationErrors(c, errors)
}

// Server errors
if err != nil {
    return response.InternalServerError(c, "Failed to process request")
}
```

## Authentication Context

When authentication middleware is implemented, extract user information:

```go
// Get authenticated user ID
userID, err := response.GetUserID(c)
if err != nil {
    return response.Unauthorized(c, "Authentication required")
}

// Use user ID in business logic
data := getDataForUser(userID)
```

## Response Patterns

### Single Resource

```go
return response.OK(c, resource)
```

### List of Resources

```go
return response.OKWithMessage(c, "Resources retrieved successfully", resources)
```

### Paginated Results

```go
return response.Paginated(c, resources, page, limit, total)
```

### Created Resource

```go
return response.Created(c, newResource)
```

### No Content

```go
return response.NoContent(c)
```

## Route Registration

Handlers are registered in the main router:

```go
func SetupRoutes(app *fiber.App, logger *config.Logger) {
    api := app.Group("/api")

    // API routes
    api.Get("/health", routes.GetHealth)
    api.Get("/resources", routes.GetResources)
    api.Post("/resources", routes.CreateResource)
    api.Get("/resources/:id", routes.GetResource)
    api.Put("/resources/:id", routes.UpdateResource)
    api.Delete("/resources/:id", routes.DeleteResource)
}
```

## Best Practices

- Keep handlers focused on HTTP concerns (parsing, validation, response formatting)
- Move business logic to separate service layers
- Use consistent naming conventions for handlers
- Validate all input data before processing
- Return appropriate HTTP status codes
- Log important events and errors
- Use the response package for all responses to maintain consistency
- Handle edge cases and error conditions gracefully
