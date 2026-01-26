package middleware

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// ValidationRule defines a validation rule for a field
type ValidationRule struct {
	Field     string
	Required  bool
	MinLength int
	MaxLength int
	Pattern   string                // For regex validation
	Validator func(value any) error // Custom validator function
}

// ValidationConfig holds validation rules for a request type
type ValidationConfig struct {
	Rules []ValidationRule
}

// ValidateRequest creates a middleware that validates request body against provided rules
// This middleware binds the request body, validates it, and stores the validated data in context
func ValidateRequest[T any](config ValidationConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req T

		// Bind request body
		if err := c.Bind().Body(&req); err != nil {
			msg := fmt.Sprintf("Failed to bind request body in validation middleware: %v", err)
			return lib.HandleServiceError(c, lib.ErrInvalidRequest, msg)
		}

		// Validate the request
		if validationErrors := validateStruct(req, config); len(validationErrors) > 0 {
			return response.SendValidationError(c, validationErrors)
		}

		// Store validated request in context for handler use
		c.Locals("validatedRequest", req)
		return c.Next()
	}
}

// validateStruct validates a struct against the provided rules
func validateStruct(data any, config ValidationConfig) []types.ValidationError {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	// Handle pointer types
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	var validationErrors []types.ValidationError

	// Validate each field according to rules
	for _, rule := range config.Rules {
		field, found := typ.FieldByName(rule.Field)
		if !found {
			continue // Skip if field doesn't exist
		}

		fieldValue := val.FieldByName(rule.Field)
		if !fieldValue.IsValid() {
			continue
		}

		// Get the actual value
		var value any
		if fieldValue.CanInterface() {
			value = fieldValue.Interface()
		}

		// Validate the field
		if err := validateField(rule, value, field.Name); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	return validationErrors
}

// validateField validates a single field against a rule
func validateField(rule ValidationRule, value any, fieldName string) *types.ValidationError {
	strValue := fmt.Sprintf("%v", value)

	// Required field validation
	if rule.Required && strings.TrimSpace(strValue) == "" {
		return &types.ValidationError{
			Field:   strings.ToLower(fieldName),
			Message: fmt.Sprintf("%s is required", fieldName),
			Value:   strValue,
		}
	}

	// Skip other validations if field is empty and not required
	if strings.TrimSpace(strValue) == "" {
		return nil
	}

	// Minimum length validation
	if rule.MinLength > 0 && len(strValue) < rule.MinLength {
		return &types.ValidationError{
			Field:   strings.ToLower(fieldName),
			Message: fmt.Sprintf("%s must be at least %d characters long", fieldName, rule.MinLength),
			Value:   strValue,
		}
	}

	// Maximum length validation
	if rule.MaxLength > 0 && len(strValue) > rule.MaxLength {
		return &types.ValidationError{
			Field:   strings.ToLower(fieldName),
			Message: fmt.Sprintf("%s must not exceed %d characters", fieldName, rule.MaxLength),
			Value:   strValue,
		}
	}

	// Custom validator function
	if rule.Validator != nil {
		if err := rule.Validator(value); err != nil {
			return &types.ValidationError{
				Field:   strings.ToLower(fieldName),
				Message: err.Error(),
				Value:   strValue,
			}
		}
	}

	return nil
}

// Common validation configurations for reuse

// AuthRequestValidation validates authentication requests
var AuthRequestValidation = ValidationConfig{
	Rules: []ValidationRule{
		{
			Field:    "Email",
			Required: true,
			Validator: func(value any) error {
				email := fmt.Sprintf("%v", value)
				if !strings.Contains(email, "@") {
					return fmt.Errorf("email must be a valid email address")
				}
				return nil
			},
		},
		{
			Field:     "Password",
			Required:  true,
			MinLength: 1,
		},
	},
}

// RegisterRequestValidation validates user registration requests
var RegisterRequestValidation = ValidationConfig{
	Rules: []ValidationRule{
		{
			Field:     "Username",
			Required:  true,
			MinLength: 3,
			MaxLength: 50,
		},
		{
			Field:    "Email",
			Required: true,
			Validator: func(value any) error {
				email := fmt.Sprintf("%v", value)
				if !strings.Contains(email, "@") {
					return fmt.Errorf("email must be a valid email address")
				}
				return nil
			},
		},
		{
			Field:     "Password",
			Required:  true,
			MinLength: 6,
			MaxLength: 128,
		},
	},
}

// FileUploadValidation validates file upload requests
var FileUploadValidation = ValidationConfig{
	Rules: []ValidationRule{
		{
			Field:    "File",
			Required: true,
		},
		{
			Field:    "SubjectID",
			Required: true,
		},
	},
}

// GetValidatedRequest retrieves the validated request from context
// This helper function provides type-safe access to validated request data
func GetValidatedRequest[T any](c fiber.Ctx) (*T, error) {
	validated := c.Locals("validatedRequest")
	if validated == nil {
		return nil, fmt.Errorf("no validated request found in context")
	}

	req, ok := validated.(T)
	if !ok {
		return nil, fmt.Errorf("request type mismatch")
	}

	return &req, nil
}

// ValidateParams creates middleware for validating URL parameters
func ValidateParams(params map[string]ValidationRule) fiber.Handler {
	return func(c fiber.Ctx) error {
		var validationErrors []types.ValidationError

		for paramName, rule := range params {
			paramValue := c.Params(paramName)

			if err := validateField(rule, paramValue, paramName); err != nil {
				validationErrors = append(validationErrors, *err)
			}
		}

		if len(validationErrors) > 0 {
			return response.SendValidationError(c, validationErrors)
		}

		return c.Next()
	}
}

// Common parameter validations

// RequiredIDParam validates that an ID parameter is present and non-empty
func RequiredIDParam(paramName string) fiber.Handler {
	return ValidateParams(map[string]ValidationRule{
		paramName: {
			Field:    paramName,
			Required: true,
		},
	})
}
