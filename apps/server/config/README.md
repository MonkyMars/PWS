# Config Package

This package handles all configuration settings for the app. It loads settings from environment variables and provides them to other parts of the application.

## What It Does

- Loads configuration from `.env` file and environment variables
- Provides default values if settings are missing
- Validates that required settings are present
- Makes configuration available throughout the app

## Main Functions

### `Load()`
Loads all configuration settings. Call this once when the app starts.

```go
cfg := config.Load()
```

### `Get()`
Gets the already loaded configuration. Use this everywhere else in your app.

```go
cfg := config.Get()
port := cfg.Port
dbHost := cfg.Database.Host
```

### Configuration Sections

**App Settings:**
- `AppName` - Name of your application
- `Environment` - "development", "production", or "staging"
- `Port` - Port to run the server on

**Database Settings:**
- `Database.Host` - Database server address
- `Database.Port` - Database port (usually 5432)
- `Database.User` - Database username
- `Database.Password` - Database password
- `Database.Name` - Database name

**Auth Settings:**
- `Auth.AccessTokenSecret` - Secret key for access tokens
- `Auth.RefreshTokenSecret` - Secret key for refresh tokens
- `Auth.AccessTokenExpiry` - How long access tokens last
- `Auth.RefreshTokenExpiry` - How long refresh tokens last

**Cache Settings:**
- `Cache.Address` - Redis server address
- `Cache.Password` - Redis password (if needed)

## How to Use

1. **In main.go:**
```go
cfg := config.Load() // Load config once at startup
```

2. **In other files:**
```go
cfg := config.Get() // Get the loaded config
serverPort := cfg.Port
dbConnection := cfg.GetDatabaseDSN()
```

3. **Check environment:**
```go
cfg := config.Get()
if cfg.IsDevelopment() {
    // Development-only code
}
if cfg.IsProduction() {
    // Production-only code
}
```

## Environment Variables

Set these in your `.env` file:

```bash
# Required
APP_NAME=PWS
PORT=8080
ACCESS_TOKEN_SECRET=your-secret-here
REFRESH_TOKEN_SECRET=your-other-secret-here

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-db-password
DB_NAME=pws_db

# Cache (Redis)
CACHE_ADDRESS=localhost:6379
CACHE_PASSWORD=your-redis-password
...
```

The package will use sensible defaults for most settings, but secrets must be provided.
