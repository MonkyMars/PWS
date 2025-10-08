// Package config provides domain-specific configuration management for the PWS application.
// This file contains modular configuration structures that separate concerns by domain,
// making the configuration more maintainable and testable.
package config

import (
	"fmt"
	"time"

	"github.com/MonkyMars/PWS/types"
)

// DomainConfigs holds all domain-specific configurations
type DomainConfigs struct {
	App      *AppConfig
	Auth     *AuthConfig
	Database *DatabaseConfig
	Server   *ServerConfig
	Cache    *CacheConfig
	Cors     *CorsConfig
	Audit    *AuditConfig
	Health   *HealthConfig
	Google   *GoogleOAuthConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string
	Environment string
	Port        string
	LogLevel    string
	FrontendURL string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	AccessTokenSecret  string
	AccessTokenExpiry  time.Duration
	RefreshTokenSecret string
	RefreshTokenExpiry time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxConns     int
	MinConns     int
	MaxLifetime  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// CacheConfig holds Redis cache configuration
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

// CorsConfig holds CORS configuration
type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

// AuditConfig holds audit logging configuration
type AuditConfig struct {
	BatchSize     int
	ChannelSize   int
	Enabled       bool
	FlushTime     time.Duration
	MaxFailures   int
	MaxRetries    int
	RetentionDays int
	RetryDelay    time.Duration
}

// HealthConfig holds health monitoring configuration
type HealthConfig struct {
	BatchSize      int
	ChannelSize    int
	Enabled        bool
	FlushTime      time.Duration
	ReportInterval time.Duration
	MaxFailures    int
	MaxRetries     int
	RetentionDays  int
	Services       []string
	RetryDelay     time.Duration
}

// GoogleOAuthConfig holds Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// LoadDomainConfigs loads all domain-specific configurations
func LoadDomainConfigs() *DomainConfigs {
	return &DomainConfigs{
		App:      loadAppConfig(),
		Auth:     loadAuthConfig(),
		Database: loadDatabaseConfig(),
		Server:   loadServerConfig(),
		Cache:    loadCacheConfig(),
		Cors:     loadCorsConfig(),
		Audit:    loadAuditConfig(),
		Health:   loadHealthConfig(),
		Google:   loadGoogleConfig(),
	}
}

// Validate validates all domain configurations
func (dc *DomainConfigs) Validate() error {
	validators := []func() error{
		dc.App.Validate,
		dc.Auth.Validate,
		dc.Database.Validate,
		dc.Server.Validate,
		dc.Cache.Validate,
		dc.Cors.Validate,
		dc.Audit.Validate,
		dc.Health.Validate,
		dc.Google.Validate,
	}

	for _, validate := range validators {
		if err := validate(); err != nil {
			return err
		}
	}

	return nil
}

// ToLegacyConfig converts domain configs to the legacy Config struct for backward compatibility
func (dc *DomainConfigs) ToLegacyConfig() *Config {
	return &Config{
		AppName:     dc.App.Name,
		Environment: dc.App.Environment,
		Port:        dc.App.Port,
		LogLevel:    dc.App.LogLevel,
		FrontendURL: dc.App.FrontendURL,
		Auth: types.AuthConfig{
			AccessTokenSecret:  dc.Auth.AccessTokenSecret,
			AccessTokenExpiry:  dc.Auth.AccessTokenExpiry,
			RefreshTokenSecret: dc.Auth.RefreshTokenSecret,
			RefreshTokenExpiry: dc.Auth.RefreshTokenExpiry,
		},
		Google: types.GoogleConfig{
			ClientID:     dc.Google.ClientID,
			ClientSecret: dc.Google.ClientSecret,
			RedirectURL:  dc.Google.RedirectURL,
		},
		Database: types.DatabaseConfig{
			Host:         dc.Database.Host,
			Port:         dc.Database.Port,
			User:         dc.Database.User,
			Password:     dc.Database.Password,
			Name:         dc.Database.Name,
			SSLMode:      dc.Database.SSLMode,
			MaxConns:     dc.Database.MaxConns,
			MinConns:     dc.Database.MinConns,
			MaxLifetime:  dc.Database.MaxLifetime,
			ReadTimeout:  dc.Database.ReadTimeout,
			WriteTimeout: dc.Database.WriteTimeout,
		},
		Server: types.ServerConfig{
			ReadTimeout:  dc.Server.ReadTimeout,
			WriteTimeout: dc.Server.WriteTimeout,
			IdleTimeout:  dc.Server.IdleTimeout,
		},
		Cache: types.CacheConfig{
			Address:         dc.Cache.Address,
			Username:        dc.Cache.Username,
			Password:        dc.Cache.Password,
			DB:              dc.Cache.DB,
			PoolSize:        dc.Cache.PoolSize,
			MinIdleConns:    dc.Cache.MinIdleConns,
			MaxIdleConns:    dc.Cache.MaxIdleConns,
			PoolTimeout:     dc.Cache.PoolTimeout,
			IdleTimeout:     dc.Cache.IdleTimeout,
			DialTimeout:     dc.Cache.DialTimeout,
			ReadTimeout:     dc.Cache.ReadTimeout,
			WriteTimeout:    dc.Cache.WriteTimeout,
			MaxRetries:      dc.Cache.MaxRetries,
			MinRetryBackoff: dc.Cache.MinRetryBackoff,
			MaxRetryBackoff: dc.Cache.MaxRetryBackoff,
		},
		Cors: types.CorsConfig{
			AllowOrigins:     dc.Cors.AllowOrigins,
			AllowMethods:     dc.Cors.AllowMethods,
			AllowHeaders:     dc.Cors.AllowHeaders,
			AllowCredentials: dc.Cors.AllowCredentials,
		},
		Audit: types.AuditConfig{
			BatchSize:     dc.Audit.BatchSize,
			ChannelSize:   dc.Audit.ChannelSize,
			Enabled:       dc.Audit.Enabled,
			FlushTime:     dc.Audit.FlushTime,
			MaxFailures:   dc.Audit.MaxFailures,
			MaxRetries:    dc.Audit.MaxRetries,
			RetentionDays: dc.Audit.RetentionDays,
			RetryDelay:    dc.Audit.RetryDelay,
		},
		Health: types.HealthConfig{
			BatchSize:      dc.Health.BatchSize,
			ChannelSize:    dc.Health.ChannelSize,
			Enabled:        dc.Health.Enabled,
			FlushTime:      dc.Health.FlushTime,
			ReportInterval: dc.Health.ReportInterval,
			MaxFailures:    dc.Health.MaxFailures,
			MaxRetries:     dc.Health.MaxRetries,
			RetentionDays:  dc.Health.RetentionDays,
			Services:       dc.Health.Services,
			RetryDelay:     dc.Health.RetryDelay,
		},
	}
}

// Domain-specific loaders
func loadAppConfig() *AppConfig {
	return &AppConfig{
		Name:        getEnv("APP_NAME", "PWS"),
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8082"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		FrontendURL: getEnv("FRONTEND_URL", ""),
	}
}

func loadAuthConfig() *AuthConfig {
	return &AuthConfig{
		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET", ""),
		AccessTokenExpiry:  getEnvDuration("ACCESS_TOKEN_EXPIRY", 15*time.Minute),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", ""),
		RefreshTokenExpiry: getEnvDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
	}
}

func loadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnvInt("DB_PORT", 5432),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", ""),
		Name:         getEnv("DB_NAME", "postgres"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		MaxConns:     getEnvInt("DB_MAX_CONNS", 25),
		MinConns:     getEnvInt("DB_MIN_CONNS", 5),
		MaxLifetime:  getEnvDuration("DB_MAX_LIFETIME", 1*time.Hour),
		ReadTimeout:  getEnvDuration("DB_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getEnvDuration("DB_WRITE_TIMEOUT", 30*time.Second),
	}
}

func loadServerConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
		WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
	}
}

func loadCacheConfig() *CacheConfig {
	return &CacheConfig{
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
	}
}

func loadCorsConfig() *CorsConfig {
	return &CorsConfig{
		AllowOrigins:     getEnvSlice("CORS_ALLOW_ORIGINS", []string{"http://localhost:5173", "http://localhost:3000"}),
		AllowMethods:     getEnvSlice("CORS_ALLOW_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		AllowHeaders:     getEnvSlice("CORS_ALLOW_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		AllowCredentials: getEnvBool("CORS_ALLOW_CREDENTIALS", true),
	}
}

func loadAuditConfig() *AuditConfig {
	return &AuditConfig{
		BatchSize:     getEnvInt("AUDIT_BATCH_SIZE", 50),
		ChannelSize:   getEnvInt("AUDIT_CHANNEL_SIZE", 1000),
		Enabled:       getEnvBool("AUDIT_ENABLED", true),
		FlushTime:     getEnvDuration("AUDIT_FLUSH_TIME", 30*time.Second),
		MaxFailures:   getEnvInt("AUDIT_MAX_FAILURES", 10),
		MaxRetries:    getEnvInt("AUDIT_MAX_RETRIES", 3),
		RetentionDays: getEnvInt("AUDIT_RETENTION_DAYS", 90),
		RetryDelay:    getEnvDuration("AUDIT_RETRY_DELAY", 3*time.Second),
	}
}

func loadHealthConfig() *HealthConfig {
	return &HealthConfig{
		BatchSize:      getEnvInt("HEALTH_BATCH_SIZE", 50),
		ChannelSize:    getEnvInt("HEALTH_CHANNEL_SIZE", 1000),
		Enabled:        getEnvBool("HEALTH_ENABLED", false),
		FlushTime:      getEnvDuration("HEALTH_FLUSH_INTERVAL", 15*time.Minute),
		ReportInterval: getEnvDuration("HEALTH_REPORT_INTERVAL", 5*time.Minute),
		MaxFailures:    getEnvInt("HEALTH_MAX_FAILURES", 10),
		MaxRetries:     getEnvInt("HEALTH_MAX_RETRIES", 3),
		RetentionDays:  getEnvInt("HEALTH_RETENTION_DAYS", 21),
		RetryDelay:     getEnvDuration("HEALTH_RETRY_DELAY", 1*time.Minute),
	}
}

func loadGoogleConfig() *GoogleOAuthConfig {
	return &GoogleOAuthConfig{
		ClientID:     getEnv("GOOGLE_OAUTH_CLIENT_ID", ""),
		ClientSecret: getEnv("GOOGLE_OAUTH_CLIENT_SECRET", ""),
		RedirectURL:  getEnv("GOOGLE_OAUTH_REDIRECT_URL", ""),
	}
}

// Domain-specific validation methods
func (ac *AppConfig) Validate() error {
	if ac.Name == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if ac.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if ac.Environment != "development" && ac.Environment != "production" && ac.Environment != "staging" {
		return fmt.Errorf("ENVIRONMENT must be one of: development, production, staging")
	}
	return nil
}

func (ac *AuthConfig) Validate() error {
	if ac.AccessTokenSecret == "" {
		return fmt.Errorf("ACCESS_TOKEN_SECRET is required")
	}
	if ac.RefreshTokenSecret == "" {
		return fmt.Errorf("REFRESH_TOKEN_SECRET is required")
	}

	// Environment-specific validation
	if getEnv("ENVIRONMENT", "development") == "production" {
		if len(ac.AccessTokenSecret) < 32 {
			return fmt.Errorf("ACCESS_TOKEN_SECRET must be at least 32 characters in production")
		}
		if len(ac.RefreshTokenSecret) < 32 {
			return fmt.Errorf("REFRESH_TOKEN_SECRET must be at least 32 characters in production")
		}
	} else {
		if len(ac.AccessTokenSecret) < 16 {
			return fmt.Errorf("ACCESS_TOKEN_SECRET must be at least 16 characters")
		}
		if len(ac.RefreshTokenSecret) < 16 {
			return fmt.Errorf("REFRESH_TOKEN_SECRET must be at least 16 characters")
		}
	}
	return nil
}

func (dc *DatabaseConfig) Validate() error {
	if dc.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if dc.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if dc.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if dc.MaxConns < 1 {
		return fmt.Errorf("DB_MAX_CONNS must be at least 1")
	}
	if dc.MinConns < 0 {
		return fmt.Errorf("DB_MIN_CONNS cannot be negative")
	}
	if dc.MinConns > dc.MaxConns {
		return fmt.Errorf("DB_MIN_CONNS cannot be greater than DB_MAX_CONNS")
	}
	return nil
}

func (sc *ServerConfig) Validate() error {
	if sc.ReadTimeout <= 0 {
		return fmt.Errorf("SERVER_READ_TIMEOUT must be positive")
	}
	if sc.WriteTimeout <= 0 {
		return fmt.Errorf("SERVER_WRITE_TIMEOUT must be positive")
	}
	if sc.IdleTimeout <= 0 {
		return fmt.Errorf("SERVER_IDLE_TIMEOUT must be positive")
	}
	return nil
}

func (cc *CacheConfig) Validate() error {
	if cc.PoolSize < 1 {
		return fmt.Errorf("CACHE_POOL_SIZE must be at least 1")
	}
	if cc.MinIdleConns < 0 {
		return fmt.Errorf("CACHE_MIN_IDLE_CONNS cannot be negative")
	}
	if cc.MaxIdleConns < cc.MinIdleConns {
		return fmt.Errorf("CACHE_MAX_IDLE_CONNS cannot be less than CACHE_MIN_IDLE_CONNS")
	}
	return nil
}

func (cc *CorsConfig) Validate() error {
	if len(cc.AllowOrigins) == 0 {
		return fmt.Errorf("CORS_ALLOW_ORIGINS cannot be empty")
	}
	if len(cc.AllowMethods) == 0 {
		return fmt.Errorf("CORS_ALLOW_METHODS cannot be empty")
	}
	return nil
}

func (ac *AuditConfig) Validate() error {
	if ac.Enabled {
		if ac.BatchSize <= 0 {
			return fmt.Errorf("AUDIT_BATCH_SIZE must be positive when audit is enabled")
		}
		if ac.ChannelSize <= 0 {
			return fmt.Errorf("AUDIT_CHANNEL_SIZE must be positive when audit is enabled")
		}
		if ac.FlushTime <= 0 {
			return fmt.Errorf("AUDIT_FLUSH_TIME must be positive when audit is enabled")
		}
	}
	return nil
}

func (hc *HealthConfig) Validate() error {
	if hc.Enabled {
		if hc.BatchSize <= 0 {
			return fmt.Errorf("HEALTH_BATCH_SIZE must be positive when health monitoring is enabled")
		}
		if hc.ChannelSize <= 0 {
			return fmt.Errorf("HEALTH_CHANNEL_SIZE must be positive when health monitoring is enabled")
		}
		if hc.FlushTime <= 0 {
			return fmt.Errorf("HEALTH_FLUSH_TIME must be positive when health monitoring is enabled")
		}
		if hc.ReportInterval <= 0 {
			return fmt.Errorf("HEALTH_REPORT_INTERVAL must be positive when health monitoring is enabled")
		}
	}
	return nil
}

func (gc *GoogleOAuthConfig) Validate() error {
	// Only validate if any Google OAuth field is set
	if gc.ClientID != "" || gc.ClientSecret != "" || gc.RedirectURL != "" {
		if gc.ClientID == "" {
			return fmt.Errorf("GOOGLE_OAUTH_CLIENT_ID is required when Google OAuth is configured")
		}
		if gc.ClientSecret == "" {
			return fmt.Errorf("GOOGLE_OAUTH_CLIENT_SECRET is required when Google OAuth is configured")
		}
		if gc.RedirectURL == "" {
			return fmt.Errorf("GOOGLE_OAUTH_REDIRECT_URL is required when Google OAuth is configured")
		}
	}
	return nil
}

// Helper methods for domain configs
func (ac *AppConfig) IsProduction() bool {
	return ac.Environment == "production"
}

func (ac *AppConfig) IsDevelopment() bool {
	return ac.Environment == "development"
}

func (dc *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dc.User, dc.Password, dc.Host, dc.Port, dc.Name, dc.SSLMode)
}

func (ac *AppConfig) GetServerAddress() string {
	return ":" + ac.Port
}
