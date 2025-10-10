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

func NewCookieService() *CookieService {
	return &CookieService{}
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
	cfg := config.Get()

	// Only use Secure cookies in production
	isSecure := cfg.IsProduction()
	// Use Lax SameSite in development for better compatibility
	sameSite := "Lax"
	if cfg.IsProduction() {
		sameSite = "Strict"
	}

	accessCookie := fiber.Cookie{
		Name:     lib.AccessTokenCookieName,
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: sameSite,
		Expires:  time.Now().Add(co.AccessTokenExpiry),
	}
	refreshCookie := fiber.Cookie{
		Name:     lib.RefreshTokenCookieName,
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: sameSite,
		Expires:  time.Now().Add(co.RefreshTokenExpiry),
	}
	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)
}

func (co *CookieService) ClearAuthCookies(c fiber.Ctx) {
	cfg := config.Get()

	// Only use Secure cookies in production
	isSecure := cfg.IsProduction()
	// Use Lax SameSite in development for better compatibility
	sameSite := "Lax"
	if cfg.IsProduction() {
		sameSite = "Strict"
	}

	accessCookie := fiber.Cookie{
		Name:     lib.AccessTokenCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: sameSite,
		Expires:  time.Now().Add(-time.Hour),
	}
	refreshCookie := fiber.Cookie{
		Name:     lib.RefreshTokenCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: sameSite,
		Expires:  time.Now().Add(-time.Hour),
	}
	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)
}

// CookieServiceInterface defines the methods for cookie management
type CookieServiceInterface interface {
	SetAuthCookies(c fiber.Ctx, accessToken, refreshToken string)
	ClearAuthCookies(c fiber.Ctx)
	GetCookieOptions() *CookieService
}
