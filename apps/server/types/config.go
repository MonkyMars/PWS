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

type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}
