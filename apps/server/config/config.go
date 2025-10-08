// Package config provides centralized configuration management for the PWS application.
// This package combines all configuration logic from API and database configurations
// into a single, easily accessible system with proper validation and defaults.
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/types"
)

// Config holds all application configuration values loaded from environment variables.
// This struct centralizes all configuration management and provides type-safe access
// to both API and database configuration values with appropriate defaults.
//
// Note: This is maintained for backward compatibility. New code should use DomainConfigs.
type Config struct {
	// Application Settings
	AppName     string
	Environment string
	Port        string
	LogLevel    string
	FrontendURL string

	// Auth Settings
	Auth types.AuthConfig

	// Google OAuth Settings
	Google types.GoogleConfig

	// Database Settings
	Database types.DatabaseConfig

	// Server Settings
	Server types.ServerConfig

	// Cache Settings
	Cache types.CacheConfig

	// CORS Settings
	Cors types.CorsConfig

	// Audit Settings
	Audit types.AuditConfig

	// Health Check Settings
	Health types.HealthConfig

	// Domain configs for better organization
	domains *DomainConfigs
}

var (
	configInstance   *Config
	domainConfigs    *DomainConfigs
	configOnce       sync.Once
	domainConfigOnce sync.Once
)

// Load loads the configuration only once using singleton pattern.
// This function ensures that configuration is loaded exactly once during the
// application lifecycle, improving performance and consistency.
//
// Returns a pointer to the loaded Config struct containing all application settings.
func Load() *Config {
	configOnce.Do(func() {
		// Load domain configs first
		domainConfigs = LoadDomainConfigs()

		// Validate domain configs
		if err := domainConfigs.Validate(); err != nil {
			log.Fatalf("Domain configuration validation failed: %v", err)
		}

		// Convert to legacy config for backward compatibility
		configInstance = domainConfigs.ToLegacyConfig()
		configInstance.domains = domainConfigs

		// Validate legacy config for additional checks
		if err := configInstance.Validate(); err != nil {
			log.Fatalf("Configuration validation failed: %v", err)
		}
	})
	return configInstance
}

// LoadDomains loads domain-specific configurations
func LoadDomains() *DomainConfigs {
	domainConfigOnce.Do(func() {
		domainConfigs = LoadDomainConfigs()
		if err := domainConfigs.Validate(); err != nil {
			log.Fatalf("Domain configuration validation failed: %v", err)
		}
	})
	return domainConfigs
}

// GetDomains returns the domain configurations
func GetDomains() *DomainConfigs {
	if domainConfigs == nil {
		return LoadDomains()
	}
	return domainConfigs
}

// Get returns the already loaded configuration instance.
// This function provides access to the singleton configuration instance.
// It panics if Load() has not been called first, ensuring proper initialization order.
//
// Returns a pointer to the loaded Config struct.
// Panics if configuration has not been loaded via Load().
func Get() *Config {
	if configInstance == nil {
		panic("Configuration not loaded. Call config.Load() first.")
	}
	return configInstance
}

// Validate checks if the configuration is valid and all required fields are set
func (c *Config) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("APP_NAME is required")
	}

	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	// Validate database configuration
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// Validate auth secrets - required in all environments
	if c.Auth.AccessTokenSecret == "" {
		return fmt.Errorf("ACCESS_TOKEN_SECRET is required")
	}
	if c.Auth.RefreshTokenSecret == "" {
		return fmt.Errorf("REFRESH_TOKEN_SECRET is required")
	}

	// Additional validation in production
	if c.IsProduction() {
		if len(c.Auth.AccessTokenSecret) < 32 {
			return fmt.Errorf("ACCESS_TOKEN_SECRET must be at least 32 characters in production")
		}
		if len(c.Auth.RefreshTokenSecret) < 32 {
			return fmt.Errorf("REFRESH_TOKEN_SECRET must be at least 32 characters in production")
		}
	} else {
		// Minimum length in non-production environments
		if len(c.Auth.AccessTokenSecret) < 16 {
			return fmt.Errorf("ACCESS_TOKEN_SECRET must be at least 16 characters")
		}
		if len(c.Auth.RefreshTokenSecret) < 16 {
			return fmt.Errorf("REFRESH_TOKEN_SECRET must be at least 16 characters")
		}
	}

	// Validate cache configuration
	if c.Cache.PoolSize < 1 {
		return fmt.Errorf("CACHE_POOL_SIZE must be at least 1")
	}
	if c.Cache.MinIdleConns < 0 {
		return fmt.Errorf("CACHE_MIN_IDLE_CONNS cannot be negative")
	}
	if c.Cache.MaxIdleConns < c.Cache.MinIdleConns {
		return fmt.Errorf("CACHE_MAX_IDLE_CONNS cannot be less than CACHE_MIN_IDLE_CONNS")
	}

	// Validate pool settings
	if c.Database.MaxConns < 1 {
		return fmt.Errorf("DB_MAX_CONNS must be at least 1")
	}
	if c.Database.MinConns < 0 {
		return fmt.Errorf("DB_MIN_CONNS cannot be negative")
	}
	if c.Database.MinConns > c.Database.MaxConns {
		return fmt.Errorf("DB_MIN_CONNS cannot be greater than DB_MAX_CONNS")
	}

	// Validate auth settings

	// Validate cache settings
	if c.Cache.Address == "" {
		log.Println("Warning: CACHE_ADDRESS is not set, caching will be disabled")
	}

	// Validate Google OAuth settings (optional - only if any Google OAuth field is set)
	if c.Google.ClientID != "" || c.Google.ClientSecret != "" || c.Google.RedirectURL != "" {
		if c.Google.ClientID == "" {
			return fmt.Errorf("GOOGLE_OAUTH_CLIENT_ID is required when Google OAuth is configured")
		}
		if c.Google.ClientSecret == "" {
			return fmt.Errorf("GOOGLE_OAUTH_CLIENT_SECRET is required when Google OAuth is configured")
		}
		if c.Google.RedirectURL == "" {
			return fmt.Errorf("GOOGLE_OAUTH_REDIRECT_URL is required when Google OAuth is configured")
		}
	}

	if c.Environment != "development" && c.Environment != "production" && c.Environment != "staging" {
		return fmt.Errorf("ENVIRONMENT must be one of: development, production, staging")
	}

	// All validations passed
	return nil
}

// IsProduction returns true if the application is running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if the application is running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetDatabaseDSN returns a formatted database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetServerAddress returns the formatted server address
func (c *Config) GetServerAddress() string {
	return ":" + c.Port
}

// PrintConfig prints the current configuration (excluding sensitive data)
func (c *Config) PrintConfig() {
	// Config printing disabled to reduce log noise
	// Enable specific log lines below if needed for debugging
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if defaultValue != "" {
		log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Invalid duration value for %s: %s, using default: %v", key, value, defaultValue)
	}
	return defaultValue
}

// getEnvBool retrieves a boolean environment variable or returns the default value if not set or invalid.
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Invalid boolean value for %s: %s, using default: %v", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return defaultValue
}

func ValidateConfig() bool {
	cfg := Load()
	if err := cfg.Validate(); err != nil {
		log.Printf("Configuration validation error: %v", err)
		return false
	}
	return true
}

// ValidateDomainConfigs validates domain-specific configurations
func ValidateDomainConfigs() bool {
	domains := LoadDomains()
	if err := domains.Validate(); err != nil {
		log.Printf("Domain configuration validation error: %v", err)
		return false
	}
	return true
}

// GetAppConfig returns the application configuration domain
func GetAppConfig() *AppConfig {
	return GetDomains().App
}

// GetAuthConfig returns the authentication configuration domain
func GetAuthConfig() *AuthConfig {
	return GetDomains().Auth
}

// GetDatabaseConfig returns the database configuration domain
func GetDatabaseConfig() *DatabaseConfig {
	return GetDomains().Database
}

// GetServerConfig returns the server configuration domain
func GetServerConfig() *ServerConfig {
	return GetDomains().Server
}

// GetCacheConfig returns the cache configuration domain
func GetCacheConfig() *CacheConfig {
	return GetDomains().Cache
}

// GetCorsConfig returns the CORS configuration domain
func GetCorsConfig() *CorsConfig {
	return GetDomains().Cors
}

// GetAuditConfig returns the audit configuration domain
func GetAuditConfig() *AuditConfig {
	return GetDomains().Audit
}

// GetHealthConfig returns the health configuration domain
func GetHealthConfig() *HealthConfig {
	return GetDomains().Health
}

// GetGoogleConfig returns the Google OAuth configuration domain
func GetGoogleConfig() *types.GoogleConfig {
	google := GetDomains().Google
	return &types.GoogleConfig{
		ClientID:     google.ClientID,
		ClientSecret: google.ClientSecret,
		RedirectURL:  google.RedirectURL,
	}
}
