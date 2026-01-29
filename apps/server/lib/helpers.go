package lib

import (
	"strconv"
	"strings"

	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func GetUserFromContext(c fiber.Ctx) *types.User {
	claimsInterface := c.Locals("claims")
	if claimsInterface == nil {
		return nil
	}

	claims, ok := claimsInterface.(*types.AuthClaims)
	if !ok || claims == nil {
		return nil
	}

	return &types.User{
		Id:    claims.Sub,
		Email: claims.Email,
		Role:  claims.Role,
	}
}

func HasPrivileges(c fiber.Ctx) bool {
	user := GetUserFromContext(c)
	if user == nil {
		return false
	}
	return user.Role == RoleAdmin || user.Role == RoleTeacher
}

func GetQueryParamAsInt(c fiber.Ctx, key string, defaultValue, maxValue int) int {
	valueStr := c.Query(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil || value <= 0 {
		return defaultValue
	}

	if value > maxValue {
		return maxValue
	}

	return value
}

func GetParams(c fiber.Ctx, keys map[string]bool) (map[string]string, error) {
	params := make(map[string]string)

	for key, required := range keys {
		val := c.Params(key)
		if val == "" {
			if required {
				return nil, ErrMissingParameter
			}
			continue
		}
		params[key] = val
	}

	return params, nil
}

func GetQueryParams(c fiber.Ctx, keys map[string]bool) (map[string]string, error) {
	params := make(map[string]string)

	for key, required := range keys {
		val := c.Query(key)
		if val == "" {
			if required {
				return nil, ErrMissingParameter
			}
			continue
		}
		params[key] = val
	}

	return params, nil
}

func ValidatePasswordStrength(password string) error {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:',.<>?/", char):
			hasSpecial = true
		}
	}
	if hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial {
		return nil
	}
	return ErrWeakPassword
}
