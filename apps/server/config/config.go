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

	// API Settings
	AnonKey string

	// Database Settings
	Database DatabaseConfig

	// Server Settings
	Server ServerConfig
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
		log.Println("Loading centralized configuration...")
		configInstance = &Config{
			// Application Settings
			AppName:     getEnv("APP_NAME", "PWS"),
			Environment: getEnv("ENVIRONMENT", "development"),
			Port:        getEnv("PORT", "8082"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),

			// API Settings
			AnonKey: getEnv("ANON_KEY", ""),

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
		}

		// Validate configuration
		if err := configInstance.Validate(); err != nil {
			log.Fatalf("Configuration validation failed: %v", err)
		}

		log.Println("Centralized configuration loaded successfully")
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
	log.Printf("=== Configuration ===")
	log.Printf("App Name: %s", c.AppName)
	log.Printf("Environment: %s", c.Environment)
	log.Printf("Port: %s", c.Port)
	log.Printf("Log Level: %s", c.LogLevel)
	log.Printf("Database Host: %s:%d", c.Database.Host, c.Database.Port)
	log.Printf("Database Name: %s", c.Database.Name)
	log.Printf("Database User: %s", c.Database.User)
	log.Printf("Database SSL Mode: %s", c.Database.SSLMode)
	log.Printf("Database Max Connections: %d", c.Database.MaxConns)
	log.Printf("Database Min Connections: %d", c.Database.MinConns)
	log.Printf("Database Max Idle Time: %v", c.Database.MaxIdleTime)
	log.Printf("Database Max Lifetime: %v", c.Database.MaxLifetime)
	log.Printf("Server Read Timeout: %v", c.Server.ReadTimeout)
	log.Printf("Server Write Timeout: %v", c.Server.WriteTimeout)
	log.Printf("=====================")
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
