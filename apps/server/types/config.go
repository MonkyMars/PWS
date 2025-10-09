package types

import "time"

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
	MaxLifetime  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
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

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

type AuditConfig struct {
	BatchSize     int           `json:"batch_size"`
	FlushTime     time.Duration `json:"flush_time"`
	ChannelSize   int           `json:"channel_size"`
	MaxRetries    int           `json:"max_retries"`
	MaxFailures   int           `json:"max_failures"`
	RetentionDays int           `json:"retention_days"`
	Enabled       bool          `json:"enabled"`
	RetryDelay    time.Duration `json:"retry_delay"`
}

type HealthConfig struct {
	BatchSize      int           `json:"batch_size"`
	FlushTime      time.Duration `json:"flush_time"`
	ReportInterval time.Duration `json:"report_interval"`
	ChannelSize    int           `json:"channel_size"`
	MaxRetries     int           `json:"max_retries"`
	MaxFailures    int           `json:"max_failures"`
	RetentionDays  int           `json:"retention_days"`
	Enabled        bool          `json:"enabled"`
	Services       []string      `json:"services"`
	RetryDelay     time.Duration `json:"retry_delay"`
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
