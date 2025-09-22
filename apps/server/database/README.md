# Database Package

This package handles all database connections and operations. It provides a centralized way to connect to PostgreSQL and perform database queries throughout the application.

## What It Does

- Establishes and manages PostgreSQL database connections
- Provides connection pooling for better performance
- Offers utility functions for common database operations
- Handles database initialization and cleanup
- Uses go-pg as the PostgreSQL driver

## Main Files

### `database.go`
Contains the core database connection logic.

**Functions:**

**`Connect()`** - Establishes database connection
```go
// Creates a new database connection using config settings
// Returns a DB instance that wraps go-pg
func Connect() (*DB, error)
```

**`Initialize()`** - Sets up the database singleton
```go
// Initializes the global database instance
// Call this once at application startup
func Initialize() error
```

**`GetInstance()`** - Gets the database connection
```go
// Returns the global database instance
// Use this throughout your application
func GetInstance() *DB
```

**`Close()`** - Closes database connection
```go
// Closes the database connection
// Call this during application shutdown
func Close() error
```

### `query.go`
Contains common query utilities and helpers.

### `actions.go`
Contains specific database operations and queries.

## How to Use

### Setting Up Database Connection

1. **In main.go:**
```go
// Initialize database connection
err := database.Initialize()
if err != nil {
    log.Fatalf("Database initialization failed: %v", err)
}

// Close connection on shutdown
defer database.Close()
```

2. **In your services:**
```go
// Get database instance
db := database.GetInstance()

// Use for queries
var users []types.User
err := db.Model(&users).Select()
```

### Database Configuration

The database package uses configuration from the config package:

```bash
# Environment variables
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=pws_db
DB_SSLMODE=disable

# Or use a connection string
DB_CONNECTION_STRING=postgres://user:pass@localhost:5432/dbname?sslmode=disable
```

### Connection Pooling

The database connection includes pooling settings:

- **MaxConns**: Maximum number of connections in pool
- **MinConns**: Minimum idle connections to maintain
- **MaxLifetime**: Maximum time a connection can be reused
- **ReadTimeout**: Timeout for read operations
- **WriteTimeout**: Timeout for write operations

## go-pg Usage Examples

### Basic Queries

**Select records:**
```go
db := database.GetInstance()

// Select all users
var users []types.User
err := db.Model(&users).Select()

// Select user by ID
var user types.User
err := db.Model(&user).Where("id = ?", userID).Select()

// Select with conditions
var users []types.User
err := db.Model(&users).Where("role = ?", "admin").Select()
```

**Insert records:**
```go
user := &types.User{
    Id:       uuid.New(),
    Username: "newuser",
    Email:    "user@example.com",
    Role:     "user",
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

_, err := db.Model(user).Insert()
```

**Update records:**
```go
// Update specific fields
_, err := db.Model(&user).
    Set("username = ?", newUsername).
    Set("updated_at = ?", time.Now()).
    Where("id = ?", userID).
    Update()

// Update entire model
user.Username = "updated"
user.UpdatedAt = time.Now()
_, err := db.Model(&user).Where("id = ?", user.Id).Update()
```

**Delete records:**
```go
_, err := db.Model(&types.User{}).Where("id = ?", userID).Delete()
```

### Advanced Queries

**Pagination:**
```go
var users []types.User
count, err := db.Model(&users).
    Limit(limit).
    Offset(offset).
    SelectAndCount()
```

**Joins:**
```go
var users []types.User
err := db.Model(&users).
    Relation("Posts").
    Select()
```

**Transactions:**
```go
err := db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
    // Multiple operations in transaction
    _, err := tx.Model(&user).Insert()
    if err != nil {
        return err
    }
    
    _, err = tx.Model(&profile).Insert()
    return err
})
```

## Database Schema

The application expects these database tables:

**users table:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Example migration:**
```sql
-- Add indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

## Error Handling

Database operations return go-pg specific errors:

```go
import "github.com/go-pg/pg/v10"

user := &types.User{}
err := db.Model(user).Where("email = ?", email).Select()
if err != nil {
    if err == pg.ErrNoRows {
        // No user found
        return nil, errors.New("user not found")
    }
    // Other database error
    return nil, err
}
```

## Health Checks

Test database connectivity:

```go
// Simple ping
err := db.Ping(context.Background())
if err != nil {
    log.Printf("Database ping failed: %v", err)
}

// Query-based health check
var result int
_, err := db.QueryOne(pg.Scan(&result), "SELECT 1")
```

## Connection Management

**Singleton Pattern:**
The database package uses a singleton pattern to ensure only one database connection pool exists:

```go
var instance *DB

func Initialize() error {
    if instance != nil {
        return errors.New("database already initialized")
    }
    
    db, err := Connect()
    if err != nil {
        return err
    }
    
    instance = db
    return nil
}
```

**Graceful Shutdown:**
```go
func Close() error {
    if instance != nil {
        return instance.Close()
    }
    return nil
}
```

## Performance Tips

1. **Use connection pooling** - Set appropriate pool sizes
2. **Add database indexes** - On frequently queried columns
3. **Use transactions** - For multiple related operations
4. **Avoid N+1 queries** - Use joins or eager loading
5. **Monitor slow queries** - Enable query logging in development

## Best Practices

1. **Always handle errors** from database operations
2. **Use parameterized queries** to prevent SQL injection
3. **Close connections** properly during shutdown
4. **Use transactions** for data consistency
5. **Add appropriate indexes** for query performance
6. **Validate data** before database operations
7. **Use context** for query timeouts
8. **Log database errors** for debugging

## Development vs Production

**Development:**
- Enable query logging for debugging
- Use local PostgreSQL instance
- Smaller connection pools

**Production:**
- Disable verbose logging
- Use managed database services
- Larger connection pools
- Enable SSL connections
- Set appropriate timeouts

The database package provides a robust foundation for all data operations in the PWS application.