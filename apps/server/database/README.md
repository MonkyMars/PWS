# Database Package

This package provides a modular and centralized way to interact with your PostgreSQL database using go-pg. It includes connection pooling, configuration management, and proper error handling for reliable and performant database operations.

## Features

- **Connection Pooling**: Configurable connection pool with sensible defaults
- **Environment-based Configuration**: Load settings from environment variables or connection string
- **Health Monitoring**: Built-in connection health checks and pool statistics
- **Singleton Pattern**: Global database instance for consistent access across your application
- **Transaction Support**: Easy transaction management
- **Context Support**: Proper context handling for timeouts and cancellations

## Configuration

### Option 1: Environment Variables

Create a `.env` file or set environment variables:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable

# Optional pool settings (with defaults shown)
DB_MAX_CONNS=25
DB_MIN_CONNS=5
DB_MAX_IDLE_TIME=15m
DB_MAX_LIFETIME=1h
DB_READ_TIMEOUT=30s
DB_WRITE_TIMEOUT=30s
```

### Option 2: Connection String

Set a single environment variable:

```bash
DB_CONNECTION_STRING=postgres://username:password@localhost:5432/database_name?sslmode=disable
```

## Basic Usage

### 1. Initialize Database Connection

In your `main.go`:

```go
import (
    "github.com/MonkyMars/PWS/database"
    "github.com/MonkyMars/PWS/services"
)

func main() {
    // Initialize database
    err := services.InitializeDatabase()
    if err != nil {
        log.Fatalf("Database initialization failed: %v", err)
    }
    
    // Ensure cleanup on exit
    defer services.CloseDatabase()
    
    // Your application code...
}
```

### 2. Use Database in Your Code

```go
import "github.com/MonkyMars/PWS/database"

func someFunction() {
    db := database.GetInstance()
    
    // Use the database connection
    // db.Model(&user).Select()
}
```

## Repository Pattern Example

Create repository structs to organize your database operations:

```go
type UserRepository struct {
    db *database.DB
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        db: database.GetInstance(),
    }
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    _, err := r.db.ModelContext(ctx, user).Insert()
    return err
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    user := &User{}
    err := r.db.ModelContext(ctx, user).Where("id = ?", id).Select()
    return user, err
}
```

## Transactions

Use transactions for operations that need to be atomic:

```go
func (r *UserRepository) TransferOperation(ctx context.Context) error {
    return r.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
        // Multiple operations that should succeed or fail together
        _, err := tx.ModelContext(ctx, &user1).Update()
        if err != nil {
            return err
        }
        
        _, err = tx.ModelContext(ctx, &user2).Update()
        return err
    })
}
```

## Raw SQL Queries

For complex queries, you can use raw SQL:

```go
func (r *UserRepository) CustomQuery(ctx context.Context) ([]*User, error) {
    var users []*User
    
    query := `
        SELECT id, email, name 
        FROM users 
        WHERE created_at > ?
        ORDER BY name
    `
    
    _, err := r.db.QueryContext(ctx, &users, query, time.Now().AddDate(0, -1, 0))
    return users, err
}
```

## Health Monitoring

Check database health and connection pool statistics:

```go
// Check if database is healthy
err := database.GetInstance().Health()
if err != nil {
    log.Printf("Database health check failed: %v", err)
}

// Get connection pool stats
stats := database.GetInstance().GetStats()
log.Printf("Pool stats - Total: %d, Idle: %d, Hits: %d", 
    stats.TotalConns, stats.IdleConns, stats.Hits)
```

## Best Practices

1. **Always use context**: Pass context to all database operations for proper timeout handling
2. **Use repositories**: Organize your database logic using the repository pattern
3. **Handle errors properly**: Always check and handle database errors
4. **Use transactions for consistency**: Group related operations in transactions
5. **Monitor connection pool**: Keep an eye on pool statistics in production
6. **Close connections**: The package handles this automatically, but ensure proper cleanup in your main function

## Connection Pool Configuration

The connection pool is configured with sensible defaults, but you can tune it based on your needs:

- **MaxConns**: Maximum number of connections (default: 25)
- **MinConns**: Minimum idle connections (default: 5)
- **MaxIdleTime**: How long connections can be idle (default: 15 minutes)
- **MaxLifetime**: Maximum connection lifetime (default: 1 hour)
- **ReadTimeout**: Read operation timeout (default: 30 seconds)
- **WriteTimeout**: Write operation timeout (default: 30 seconds)

## Error Handling

The package wraps database errors with context. Common patterns:

```go
user, err := userRepo.GetByID(ctx, id)
if err != nil {
    if err == pg.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    return nil, fmt.Errorf("database error: %w", err)
}
```

## Migration Support

For database migrations, you can create tables programmatically:

```go
func (r *UserRepository) CreateTable(ctx context.Context) error {
    return r.db.ModelContext(ctx, (*User)(nil)).CreateTable(&orm.CreateTableOptions{
        IfNotExists: true,
    })
}
```

## Testing

For testing, you can create a separate database instance or use a test database:

```go
func setupTestDB() *database.DB {
    config := &database.Config{
        Host:     "localhost",
        Port:     5432,
        User:     "test_user",
        Password: "test_password",
        Database: "test_db",
        SSLMode:  "disable",
    }
    
    db, err := database.Connect(config)
    if err != nil {
        panic(err)
    }
    
    return db
}
```
