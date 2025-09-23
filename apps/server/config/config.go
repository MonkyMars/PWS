// Package config provides centralized configuration management for the PWS application.
// This package combines all configuration logic from API and database configurations
// into a single, easily accessible system with proper validation and defaults.
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// Config holds all application configuration values loaded from environment variables.
// This struct centralizes all configuration management and provides type-safe access
// to both API and database configuration values with appropriate defaults.
type Config struct {
	// Application Settings
	AppName     string
	Environment string
	Port        string
	LogLevel    string

	// Auth Settings
	Auth AuthConfig

	// API Settings
	Supabase SupabaseConfig

	// Google OAuth Settings
	Google GoogleConfig

	// Frontend Settings
	FrontendURL string

	// Database Settings
	Database DatabaseConfig

	// Server Settings
	Server ServerConfig

	// Cache Settings
	Cache CacheConfig
}

// DatabaseConfig holds all database-related configuration
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxConns     int
	MinConns     int
	MaxIdleTime  time.Duration
	MaxLifetime  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Alternative: Full connection string
	ConnectionString string
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

type SupabaseConfig struct {
	Url        string
	AnonKey    string
	ServiceKey string
}

type AuthConfig struct {
	AccessTokenSecret  string
	AccessTokenExpiry  time.Duration
	RefreshTokenSecret string
	RefreshTokenExpiry time.Duration
}

type CacheConfig struct {
	Address         string
	Username        string
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	MaxIdleConns    int
	PoolTimeout     time.Duration
	IdleTimeout     time.Duration
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

var (
	configInstance *Config
	configOnce     sync.Once
)

// Load loads the configuration only once using singleton pattern.
// This function ensures that configuration is loaded exactly once during the
// application lifecycle, improving performance and consistency.
//
// Returns a pointer to the loaded Config struct containing all application settings.
func Load() *Config {
	configOnce.Do(func() {
		configInstance = &Config{
			// Application Settings
			AppName:     getEnv("APP_NAME", "PWS"),
			Environment: getEnv("ENVIRONMENT", "development"),
			Port:        getEnv("PORT", "8082"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),

			// API Settings
			Supabase: SupabaseConfig{
				Url:        getEnv("SUPABASE_URL", ""),
				AnonKey:    getEnv("SUPABASE_ANON_KEY", ""),
				ServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
			},

			Auth: AuthConfig{
				AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", ""),
				AccessTokenExpiry:  getEnvDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
				RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", ""),
				RefreshTokenExpiry: getEnvDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
			},

			// Google OAuth Settings
			Google: GoogleConfig{
				ClientID:     getEnv("GOOGLE_OAUTH_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_OAUTH_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_OAUTH_REDIRECT_URL", ""),
			},

			// Frontend Settings
			FrontendURL: getEnv("FRONTEND_URL", ""),

			// Database Settings
			Database: DatabaseConfig{
				Host:             getEnv("DB_HOST", "localhost"),
				Port:             getEnvInt("DB_PORT", 5432),
				User:             getEnv("DB_USER", "postgres"),
				Password:         getEnv("DB_PASSWORD", ""),
				Name:             getEnv("DB_NAME", "postgres"),
				SSLMode:          getEnv("DB_SSLMODE", "disable"),
				MaxConns:         getEnvInt("DB_MAX_CONNS", 25),
				MinConns:         getEnvInt("DB_MIN_CONNS", 5),
				MaxIdleTime:      getEnvDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
				MaxLifetime:      getEnvDuration("DB_MAX_LIFETIME", 1*time.Hour),
				ReadTimeout:      getEnvDuration("DB_READ_TIMEOUT", 30*time.Second),
				WriteTimeout:     getEnvDuration("DB_WRITE_TIMEOUT", 30*time.Second),
				ConnectionString: getEnv("DB_CONNECTION_STRING", ""),
			},

			// Server Settings
			Server: ServerConfig{
				ReadTimeout:    getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
				WriteTimeout:   getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
				IdleTimeout:    getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
				MaxHeaderBytes: getEnvInt("SERVER_MAX_HEADER_BYTES", 1<<20), // 1MB
			},

			// Cache Settings
			Cache: CacheConfig{
				Address:         getEnv("CACHE_ADDRESS", "localhost:6379"),
				Username:        getEnv("CACHE_USERNAME", ""),
				Password:        getEnv("CACHE_PASSWORD", ""),
				DB:              getEnvInt("CACHE_DB", 0),
				PoolSize:        getEnvInt("CACHE_POOL_SIZE", 10),
				MinIdleConns:    getEnvInt("CACHE_MIN_IDLE_CONNS", 2),
				MaxIdleConns:    getEnvInt("CACHE_MAX_IDLE_CONNS", 5),
				PoolTimeout:     getEnvDuration("CACHE_POOL_TIMEOUT", 30*time.Second),
				IdleTimeout:     getEnvDuration("CACHE_IDLE_TIMEOUT", 5*time.Minute),
				DialTimeout:     getEnvDuration("CACHE_DIAL_TIMEOUT", 5*time.Second),
				ReadTimeout:     getEnvDuration("CACHE_READ_TIMEOUT", 3*time.Second),
				WriteTimeout:    getEnvDuration("CACHE_WRITE_TIMEOUT", 3*time.Second),
				MaxRetries:      getEnvInt("CACHE_MAX_RETRIES", 3),
				MinRetryBackoff: getEnvDuration("CACHE_MIN_RETRY_BACKOFF", 8*time.Millisecond),
				MaxRetryBackoff: getEnvDuration("CACHE_MAX_RETRY_BACKOFF", 512*time.Millisecond),
			},
		}

		// Validate configuration
		if err := configInstance.Validate(); err != nil {
			log.Fatalf("Configuration validation failed: %v", err)
		}
	})
	return configInstance
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

	// Validate database configuration only if no connection string is provided
	if c.Database.ConnectionString == "" {
		if c.Database.Host == "" {
			return fmt.Errorf("DB_HOST is required when DB_CONNECTION_STRING is not provided")
		}
		if c.Database.User == "" {
			return fmt.Errorf("DB_USER is required when DB_CONNECTION_STRING is not provided")
		}
		if c.Database.Name == "" {
			return fmt.Errorf("DB_NAME is required when DB_CONNECTION_STRING is not provided")
		}
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

// IsStaging returns true if the application is running in staging environment
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

// GetDatabaseDSN returns a formatted database connection string
func (c *Config) GetDatabaseDSN() string {
	if c.Database.ConnectionString != "" {
		return c.Database.ConnectionString
	}

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

func ValidateConfig() bool {
	cfg := Load()
	if err := cfg.Validate(); err != nil {
		log.Printf("Configuration validation error: %v", err)
		return false
	}
	return true
}
