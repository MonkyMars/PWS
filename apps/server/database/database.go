package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/types"
	"github.com/go-pg/pg/v10"
)

type DB struct {
	*pg.DB
}

var instance *DB

// Connect establishes a connection to the database using centralized configuration
func Connect() (*DB, error) {
	cfg := config.Get()
	dbCfg := cfg.Database

	opts := &pg.Options{
		Addr:     fmt.Sprintf("%s:%d", dbCfg.Host, dbCfg.Port),
		User:     dbCfg.User,
		Password: dbCfg.Password,
		Database: dbCfg.Name,
	}

	// Apply pool settings from configuration
	opts.PoolSize = dbCfg.MaxConns
	opts.MinIdleConns = dbCfg.MinConns
	opts.MaxConnAge = dbCfg.MaxLifetime
	opts.ReadTimeout = dbCfg.ReadTimeout
	opts.WriteTimeout = dbCfg.WriteTimeout

	db := pg.Connect(opts)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database: %s@%s:%d/%s",
		dbCfg.User, dbCfg.Host, dbCfg.Port, dbCfg.Name)

	return &DB{db}, nil
}

// Initialize sets up the global database instance using centralized configuration
func Initialize() error {
	db, err := Connect()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	instance = db
	return nil
}

// GetInstance returns the global database instance
func GetInstance() *DB {
	if instance == nil {
		log.Fatal("Database not initialized. Call Initialize() first.")
	}
	return instance
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// CloseInstance closes the global database instance
func CloseInstance() error {
	if instance != nil {
		return instance.Close()
	}
	return nil
}

// Health checks the database connection health
func (db *DB) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.Ping(ctx)
}

// GetStats returns connection pool statistics
func (db *DB) GetStats() *pg.PoolStats {
	return db.PoolStats()
}

// GetConfig returns the current database configuration
func GetConfig() types.DatabaseConfig {
	cfg := config.Get()
	return cfg.Database
}
