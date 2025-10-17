package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
	redisCtx    = context.Background()
)

// CacheService provides Redis caching functionality with connection pooling and retry logic
type CacheService struct{}

func NewCacheService() *CacheService {
	return &CacheService{}
}

// GetRedisClient returns a singleton Redis client with proper connection pooling
func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		cfg := config.Get()

		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Cache.Address,
			Username: cfg.Cache.Username,
			Password: cfg.Cache.Password,
			DB:       cfg.Cache.DB,

			// Connection pool settings
			PoolSize:        cfg.Cache.PoolSize,
			MinIdleConns:    cfg.Cache.MinIdleConns,
			MaxIdleConns:    cfg.Cache.MaxIdleConns,
			PoolTimeout:     cfg.Cache.PoolTimeout,
			ConnMaxIdleTime: cfg.Cache.IdleTimeout,

			// Timeouts
			DialTimeout:  cfg.Cache.DialTimeout,
			ReadTimeout:  cfg.Cache.ReadTimeout,
			WriteTimeout: cfg.Cache.WriteTimeout,

			// Retry settings
			MaxRetries:      cfg.Cache.MaxRetries,
			MinRetryBackoff: cfg.Cache.MinRetryBackoff,
			MaxRetryBackoff: cfg.Cache.MaxRetryBackoff,
		})
	})
	return redisClient
}

// CloseRedisConnection closes the Redis connection pool
func CloseRedisConnection() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// withRetry executes a Redis operation with exponential backoff retry logic
func (cs *CacheService) withRetry(operation func() error, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on the last attempt
		if attempt == maxRetries {
			break
		}

		// Only retry on network/connection errors, not on logical errors like key not found
		if !isRetryableError(err) {
			return err
		}

		maxBackoff := 2000 // max 2000ms = 2s
		base := 100        // 100ms base

		backoff := base * (1 << attempt) // exponential
		if backoff > maxBackoff {
			backoff = maxBackoff
		}

		// add jitter ±50%
		jitterBytes := make([]byte, 4)
		_, err = rand.Read(jitterBytes)
		if err != nil {
			// fallback to no jitter if random fails
			time.Sleep(time.Duration(backoff) * time.Millisecond)
			continue
		}
		jitter := int(uint32(jitterBytes[0])<<24 | uint32(jitterBytes[1])<<16 | uint32(jitterBytes[2])<<8 | uint32(jitterBytes[3]))
		// No need to handle negative values; uint32 avoids sign extension
		// jitter is always non-negative

		jitter = jitter % (backoff/2 + 1)
		backoffWithJitter := backoff/2 + jitter

		time.Sleep(time.Duration(backoffWithJitter) * time.Millisecond)
	}

	return fmt.Errorf("redis operation failed after %d retries: %w", maxRetries, lastErr)
}

// isRetryableError determines if an error is worth retrying
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Don't retry on nil results (key not found)
	if err == redis.Nil {
		return false
	}

	// Retry on network/connection errors
	errStr := err.Error()
	retryableErrors := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"broken pipe",
		"no such host",
		"network is unreachable",
	}

	for _, retryableErr := range retryableErrors {
		if strings.Contains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// Set sets a key with TTL and automatic retry logic
func (cs *CacheService) Set(key string, value any, ttl time.Duration) error {
	client := GetRedisClient()

	return cs.withRetry(func() error {
		return client.Set(redisCtx, key, value, ttl).Err()
	}, 3)
}

// Get retrieves a key with automatic retry logic
func (cs *CacheService) Get(key string) (string, error) {
	client := GetRedisClient()
	var result string
	var resultErr error

	err := cs.withRetry(func() error {
		val, err := client.Get(redisCtx, key).Result()
		if err == redis.Nil {
			result = ""
			resultErr = nil
			return nil // Don't retry on key not found
		}
		if err != nil {
			return err
		}
		result = val
		resultErr = nil
		return nil
	}, 3)

	if err != nil {
		return "", err
	}

	return result, resultErr
}

// Delete removes a key with automatic retry logic
func (cs *CacheService) Delete(key string) error {
	client := GetRedisClient()

	return cs.withRetry(func() error {
		return client.Del(redisCtx, key).Err()
	}, 3)
}

// Exists checks if a key exists with automatic retry logic
func (cs *CacheService) Exists(key string) (bool, error) {
	client := GetRedisClient()
	var result bool

	err := cs.withRetry(func() error {
		count, err := client.Exists(redisCtx, key).Result()
		if err != nil {
			return err
		}
		result = count > 0
		return nil
	}, 3)

	return result, err
}

// BlacklistToken adds a token's jti to the blacklist with expiration and retry logic
func (cs *CacheService) BlacklistToken(jti string, exp time.Time) error {
	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil // token already expired, no need to store
	}

	key := fmt.Sprintf("blacklist:%s", jti)
	return cs.Set(key, "true", ttl)
}

// IsTokenBlacklisted checks if a JTI exists in Redis with retry logic
func (cs *CacheService) IsTokenBlacklisted(jti uuid.UUID) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", jti.String())
	val, err := cs.Get(key)
	if err != nil {
		return false, err
	}

	return val == "true", nil
}

// SetUserSession stores user session data with TTL
func (cs *CacheService) SetUserSession(userID, sessionID string, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s:%s", userID, sessionID)
	return cs.Set(key, "active", ttl)
}

// GetUserSession retrieves user session data
func (cs *CacheService) GetUserSession(userID, sessionID string) (bool, error) {
	key := fmt.Sprintf("session:%s:%s", userID, sessionID)
	val, err := cs.Get(key)
	if err != nil {
		return false, err
	}

	return val == "active", nil
}

// DeleteUserSession removes a user session
func (cs *CacheService) DeleteUserSession(userID, sessionID string) error {
	key := fmt.Sprintf("session:%s:%s", userID, sessionID)
	return cs.Delete(key)
}

// SetRateLimit sets a rate limit counter for an IP/endpoint combination
func (cs *CacheService) SetRateLimit(ip, endpoint string, count int, ttl time.Duration) error {
	key := fmt.Sprintf("ratelimit:%s:%s", ip, endpoint)
	return cs.Set(key, count, ttl)
}

// GetRateLimit retrieves the current rate limit count for an IP/endpoint
func (cs *CacheService) GetRateLimit(ip, endpoint string) (int, error) {
	key := fmt.Sprintf("ratelimit:%s:%s", ip, endpoint)
	val, err := cs.Get(key)
	if err != nil {
		return 0, err
	}

	if val == "" {
		return 0, nil
	}

	// Simple integer parsing (you might want to use strconv.Atoi for better error handling)
	count := 0
	for _, char := range val {
		if char >= '0' && char <= '9' {
			count = count*10 + int(char-'0')
		} else {
			break
		}
	}

	return count, nil
}

// IncrementRateLimit atomically increments a rate limit counter
func (cs *CacheService) IncrementRateLimit(ip, endpoint string, ttl time.Duration) (int, error) {
	client := GetRedisClient()
	key := fmt.Sprintf("ratelimit:%s:%s", ip, endpoint)

	var result int64
	err := cs.withRetry(func() error {
		val, err := client.Incr(redisCtx, key).Result()
		if err != nil {
			return err
		}
		result = val

		// Set expiration only on first increment
		if val == 1 {
			return client.Expire(redisCtx, key, ttl).Err()
		}

		return nil
	}, 3)

	return int(result), err
}

// Ping tests the Redis connection
func (cs *CacheService) Ping() error {
	client := GetRedisClient()

	return cs.withRetry(func() error {
		return client.Ping(redisCtx).Err()
	}, 3)
}

// GetConnectionStats returns Redis connection pool statistics
func (cs *CacheService) GetConnectionStats() map[string]any {
	client := GetRedisClient()
	stats := client.PoolStats()

	return map[string]any{
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"timeouts":    stats.Timeouts,
		"total_conns": stats.TotalConns,
		"idle_conns":  stats.IdleConns,
		"stale_conns": stats.StaleConns,
	}
}

// GetRedisInfo returns Redis server information for monitoring
func (cs *CacheService) GetRedisInfo() (map[string]string, error) {
	client := GetRedisClient()
	var result map[string]string

	err := cs.withRetry(func() error {
		info, err := client.Info(redisCtx).Result()
		if err != nil {
			return err
		}

		// Parse the info string into a map
		result = make(map[string]string)
		lines := strings.Split(info, "\r\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
		return nil
	}, 3)

	return result, err
}

// TestRedisConnection performs a comprehensive Redis connection test
func (cs *CacheService) TestRedisConnection() error {
	client := GetRedisClient()

	return cs.withRetry(func() error {
		// Test basic connectivity
		if err := client.Ping(redisCtx).Err(); err != nil {
			return fmt.Errorf("ping failed: %w", err)
		}

		// Test set/get operations
		testKey := fmt.Sprintf("test:conn:%d", time.Now().UnixNano())
		testValue := "connection_test"

		if err := client.Set(redisCtx, testKey, testValue, time.Minute).Err(); err != nil {
			return fmt.Errorf("set operation failed: %w", err)
		}

		val, err := client.Get(redisCtx, testKey).Result()
		if err != nil {
			return fmt.Errorf("get operation failed: %w", err)
		}

		if val != testValue {
			return fmt.Errorf("value mismatch: expected %s, got %s", testValue, val)
		}

		// Clean up test key
		if err := client.Del(redisCtx, testKey).Err(); err != nil {
			return fmt.Errorf("delete operation failed: %w", err)
		}

		return nil
	}, 3)
}

// FlushBlacklistedTokens removes all blacklisted tokens (useful for maintenance)
func (cs *CacheService) FlushBlacklistedTokens() error {
	client := GetRedisClient()

	return cs.withRetry(func() error {
		keys, err := client.Keys(redisCtx, "blacklist:*").Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			return client.Del(redisCtx, keys...).Err()
		}

		return nil
	}, 3)
}

// GetBlacklistedTokensCount returns the number of currently blacklisted tokens
func (cs *CacheService) GetBlacklistedTokensCount() (int, error) {
	client := GetRedisClient()
	var count int

	err := cs.withRetry(func() error {
		keys, err := client.Keys(redisCtx, "blacklist:*").Result()
		if err != nil {
			return err
		}
		count = len(keys)
		return nil
	}, 3)

	return count, err
}

// GetActiveSessionsCount returns the number of active user sessions
func (cs *CacheService) GetActiveSessionsCount() (int, error) {
	client := GetRedisClient()
	var count int

	err := cs.withRetry(func() error {
		keys, err := client.Keys(redisCtx, "session:*").Result()
		if err != nil {
			return err
		}
		count = len(keys)
		return nil
	}, 3)

	return count, err
}

// GetRateLimitStatus returns current rate limit information for debugging
func (cs *CacheService) GetRateLimitStatus(ip, endpoint string) (map[string]any, error) {
	key := fmt.Sprintf("ratelimit:%s:%s", ip, endpoint)

	client := GetRedisClient()
	var result map[string]any

	err := cs.withRetry(func() error {
		// Get current count
		val, err := client.Get(redisCtx, key).Result()
		if err == redis.Nil {
			result = map[string]any{
				"count": 0,
				"ttl":   0,
			}
			return nil
		}
		if err != nil {
			return err
		}

		// Get TTL
		ttl, err := client.TTL(redisCtx, key).Result()
		if err != nil {
			return err
		}

		// Parse count
		count := 0
		for _, char := range val {
			if char >= '0' && char <= '9' {
				count = count*10 + int(char-'0')
			} else {
				break
			}
		}

		result = map[string]any{
			"count": count,
			"ttl":   int(ttl.Seconds()),
		}
		return nil
	}, 3)

	return result, err
}

type CacheServiceInterface interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)

	BlacklistToken(jti uuid.UUID, exp time.Time) error
	IsTokenBlacklisted(jti uuid.UUID) (bool, error)

	SetUserSession(userID, sessionID string, ttl time.Duration) error
	GetUserSession(userID, sessionID string) (bool, error)
	DeleteUserSession(userID, sessionID string) error

	SetRateLimit(ip, endpoint string, count int, ttl time.Duration) error
	GetRateLimit(ip, endpoint string) (int, error)
	IncrementRateLimit(ip, endpoint string, ttl time.Duration) (int, error)

	Ping() error
	GetConnectionStats() map[string]any
	GetRedisInfo() (map[string]string, error)
	TestRedisConnection() error

	FlushBlacklistedTokens() error
	GetBlacklistedTokensCount() (int, error)
	GetActiveSessionsCount() (int, error)
	GetRateLimitStatus(ip, endpoint string) (map[string]any, error)
}
