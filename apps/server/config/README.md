# Centralized Configuration System

This package provides a centralized configuration management system for the PWS application, combining all configuration logic from API, database, and server settings into one easily accessible location.

## 🎯 Overview

The centralized config system eliminates configuration duplication and provides:

- **Single Source of Truth**: All configuration in one place
- **Type Safety**: Strongly typed configuration structs
- **Environment Validation**: Automatic validation of required settings
- **Smart Defaults**: Sensible defaults for all optional settings
- **Easy Access**: Simple `config.Get()` access throughout your application
- **Development Helpers**: Built-in environment detection and debug utilities

## 🚀 Quick Start

### 1. Load Configuration (in main.go)

```go
import "github.com/MonkyMars/PWS/config"

func main() {
    // Load configuration once at startup
    cfg := config.Load()
    
    // Use throughout your application
    logger := config.SetupLogger()
    fiberConfig := config.SetupFiber()
}
```

### 2. Access Configuration Anywhere

```go
import "github.com/MonkyMars/PWS/config"

func someFunction() {
    cfg := config.Get()
    
    if cfg.IsDevelopment() {
        // Development-specific logic
    }
    
    dbDSN := cfg.GetDatabaseDSN()
    serverAddr := cfg.GetServerAddress()
}
```

## 📁 Configuration Structure

```go
type Config struct {
    // Application Settings
    AppName     string  // APP_NAME
    Environment string  // ENVIRONMENT
    Port        string  // PORT
    LogLevel    string  // LOG_LEVEL
    
    // API Settings
    AnonKey     string  // ANON_KEY
    
    // Database Settings
    Database    DatabaseConfig
    
    // Server Settings
    Server      ServerConfig
}
```

## 🔧 Environment Variables

### Required Variables
```bash
APP_NAME=PWS                    # Application name
ENVIRONMENT=development         # development, staging, production
PORT=8082                      # Server port
```

### Database Configuration

**Option 1: Individual Parameters**
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable
```

**Option 2: Connection String (overrides individual parameters)**
```bash
DB_CONNECTION_STRING=postgres://user:pass@localhost:5432/dbname?sslmode=disable
```

### Optional Settings with Defaults
```bash
# Logging
LOG_LEVEL=info                  # debug, info, warn, error

# API
ANON_KEY=                       # Your API anonymous key

# Database Pool
DB_MAX_CONNS=25                 # Maximum connections
DB_MIN_CONNS=5                  # Minimum idle connections
DB_MAX_IDLE_TIME=15m            # Max idle time
DB_MAX_LIFETIME=1h              # Max connection lifetime
DB_READ_TIMEOUT=30s             # Read timeout
DB_WRITE_TIMEOUT=30s            # Write timeout

# Server
SERVER_READ_TIMEOUT=30s         # HTTP read timeout
SERVER_WRITE_TIMEOUT=30s        # HTTP write timeout
SERVER_IDLE_TIMEOUT=120s        # HTTP idle timeout
SERVER_MAX_HEADER_BYTES=1048576 # Max header size (1MB)
```

## 🛠 Usage Examples

### Basic Application Setup

```go
package main

import (
    "github.com/MonkyMars/PWS/config"
    "github.com/MonkyMars/PWS/database"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    godotenv.Load()
    
    // Load centralized configuration
    cfg := config.Load()
    
    // Setup components using centralized config
    logger := config.SetupLogger()
    
    // Initialize database (uses config automatically)
    err := database.Initialize()
    if err != nil {
        logger.DatabaseError("initialization", err)
        panic(err)
    }
    
    logger.DatabaseConnected()
    
    // Your application logic...
}
```

### Database Usage

```go
import (
    "github.com/MonkyMars/PWS/config"
    "github.com/MonkyMars/PWS/database"
)

func setupDatabase() {
    // Database package automatically uses centralized config
    db := database.GetInstance()
    
    // Access database config if needed
    cfg := config.Get()
    if cfg.IsDevelopment() {
        cfg.PrintConfig() // Shows all config (except sensitive data)
    }
}
```

### Fiber Web Server

```go
import (
    "github.com/MonkyMars/PWS/config"
    "github.com/gofiber/fiber/v3"
)

func setupServer() *fiber.App {
    // Create Fiber app with centralized config
    app := fiber.New(config.SetupFiber())
    
    // Add centralized logging middleware
    logger := config.SetupLogger()
    app.Use(logger.HTTPMiddleware())
    
    return app
}
```

### Environment-Specific Logic

```go
func someBusinessLogic() {
    cfg := config.Get()
    
    switch {
    case cfg.IsDevelopment():
        // Development-specific behavior
        logger.Debug("Development mode active")
        
    case cfg.IsStaging():
        // Staging-specific behavior
        setupStagingFeatures()
        
    case cfg.IsProduction():
        // Production-specific behavior
        enableProductionOptimizations()
    }
}
```

## 🔍 Built-in Helper Methods

### Environment Detection
```go
cfg := config.Get()

cfg.IsDevelopment()  // true if ENVIRONMENT=development
cfg.IsStaging()      // true if ENVIRONMENT=staging
cfg.IsProduction()   // true if ENVIRONMENT=production
```

### Formatted Values
```go
cfg.GetDatabaseDSN()    // Returns formatted database connection string
cfg.GetServerAddress()  // Returns ":8082" format for server listening
```

### Configuration Display
```go
cfg.PrintConfig()  // Logs all config values (hides sensitive data)
```

## 🧪 Testing

### Test Configuration
```go
// In your tests, you can create test-specific config
func setupTestConfig() {
    os.Setenv("ENVIRONMENT", "test")
    os.Setenv("DB_NAME", "test_database")
    os.Setenv("LOG_LEVEL", "error")
    
    // Load test configuration
    cfg := config.Load()
}
```

### Mock Configuration
```go
// For unit tests, you can mock the config
func TestWithMockConfig(t *testing.T) {
    // Set test environment variables
    os.Setenv("APP_NAME", "TestApp")
    os.Setenv("ENVIRONMENT", "test")
    
    cfg := config.Load()
    assert.Equal(t, "TestApp", cfg.AppName)
}
```

## 🔒 Security Considerations

1. **Sensitive Data**: Never commit `.env` files with real credentials
2. **Production**: Use environment variables or secret management systems
3. **Validation**: Configuration validation prevents startup with invalid settings
4. **Logging**: Sensitive fields are excluded from `PrintConfig()` output

## 📝 Configuration Validation

The system automatically validates:

- ✅ Required fields are present
- ✅ Database connection parameters are valid
- ✅ Pool settings are reasonable
- ✅ Port numbers are valid
- ✅ Timeout durations are parseable

```go
// Validation happens automatically during Load()
cfg := config.Load() // Will panic if validation fails
```

## 🔄 Migration from Old Config

If you're migrating from the old `api/config` system:

1. **Replace imports**:
   ```go
   // Old
   import "github.com/MonkyMars/PWS/api/config"
   
   // New
   import "github.com/MonkyMars/PWS/config"
   ```

2. **Update function calls**:
   ```go
   // Old
   cfg := config.LoadEnv()
   logger := config.SetupLogger(cfg)
   fiberCfg := config.SetupFiber(cfg)
   
   // New
   cfg := config.Load()
   logger := config.SetupLogger()
   fiberCfg := config.SetupFiber()
   ```

3. **Update field access**:
   ```go
   // Old
   cfg.ConnectionString
   
   // New
   cfg.Database.ConnectionString
   cfg.GetDatabaseDSN() // or use helper method
   ```

## 🎛 Environment Examples

### Development (.env)
```bash
ENVIRONMENT=development
LOG_LEVEL=debug
DB_HOST=localhost
DB_USER=dev_user
DB_PASSWORD=dev_password
```

### Staging
```bash
ENVIRONMENT=staging
LOG_LEVEL=info
DB_HOST=staging-db.example.com
DB_CONNECTION_STRING=postgres://staging_user:pass@staging-db:5432/staging_db
```

### Production
```bash
ENVIRONMENT=production
LOG_LEVEL=warn
DB_CONNECTION_STRING=postgres://prod_user:secure_pass@prod-db:5432/prod_db?sslmode=require
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
```

## 🚀 Benefits

- **🎯 Centralized**: All config in one place
- **🔒 Type Safe**: Compile-time type checking
- **⚡ Performance**: Singleton pattern, loaded once
- **🛡 Validated**: Automatic validation on startup
- **🔧 Flexible**: Support for env vars and connection strings
- **📊 Observable**: Built-in logging and debugging
- **🧪 Testable**: Easy to mock and test
- **📈 Scalable**: Easy to add new configuration options

The centralized configuration system provides a robust foundation for managing all your application settings in a clean, maintainable, and scalable way.