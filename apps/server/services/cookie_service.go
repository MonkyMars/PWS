package services

import (
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/gofiber/fiber/v3"
)

type CookieService struct {
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func (co *CookieService) GetCookieOptions() *CookieService {
	cfg := config.Get()
	return &CookieService{
		AccessTokenExpiry:  cfg.Auth.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.Auth.RefreshTokenExpiry,
	}
}

func (co *CookieService) SetAuthCookies(c fiber.Ctx, accessToken, refreshToken string) {
	co = co.GetCookieOptions()
	accessCookie := fiber.Cookie{
		Name:     lib.AccessTokenCookieName,
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(co.AccessTokenExpiry),
	}
	refreshCookie := fiber.Cookie{
		Name:     lib.RefreshTokenCookieName,
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(co.RefreshTokenExpiry),
	}
	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)
}

func (co *CookieService) ClearAuthCookies(c fiber.Ctx) {
	accessCookie := fiber.Cookie{
		Name:     lib.AccessTokenCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
	}
	refreshCookie := fiber.Cookie{
		Name:     lib.RefreshTokenCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour),
	}
	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)
}
