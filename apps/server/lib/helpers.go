package lib

import (
	"strconv"

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

func GetParams(c fiber.Ctx, keys ...string) (map[string]string, error) {
	params := make(map[string]string)
	for _, key := range keys {
		val := c.Params(key)
		if val == "" {
			return nil, ErrMissingParameter
		} else {
			params[key] = val
		}
	}
	return params, nil
}
