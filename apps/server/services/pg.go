package services

import (
	"context"
	"time"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/types"
	"github.com/go-pg/pg/v10"
)

// Find executes a SELECT query to find records
// Usage: users, err := services.Find[User]("users", map[string]any{"active": true})
func Find[T any](table string, where map[string]any) (*types.QueryResult[T], error) {
	return database.Select[T](table, where)
}

// FindOne executes a SELECT query to find a single record
// Usage: user, err := services.FindOne[User]("users", map[string]any{"id": 1})
func FindOne[T any](table string, where map[string]any) (*T, error) {
	query := types.NewQuery().
		SetOperation("select").
		SetTable(table).
		SetLimit(1)

	for key, value := range where {
		query.AddWhere(key, value)
	}

	result, err := database.ExecuteQuery[T](query)
	if err != nil {
		return nil, err
	}

	if result.Single != nil {
		return result.Single, nil
	}

	return nil, nil
}

// Create executes an INSERT query
// Usage: result, err := services.Create[User]("users", map[string]any{"name": "John", "email": "john@example.com"})
func Create[T any](table string, data map[string]any) (*types.QueryResult[T], error) {
	return database.Insert[T](table, data)
}

// CreateAndReturn executes an INSERT query and returns the created record
// Usage: user, err := services.CreateAndReturn[User]("users", map[string]any{"name": "John"}, "id", "name", "created_at")
func CreateAndReturn[T any](table string, data map[string]any, returning ...string) (*T, error) {
	query := types.NewQuery().
		SetOperation("insert").
		SetTable(table).
		SetData(data).
		SetReturning(returning...)

	result, err := database.ExecuteQuery[T](query)
	if err != nil {
		return nil, err
	}

	return result.Single, nil
}

// Update executes an UPDATE query
// Usage: result, err := services.Update[User]("users", map[string]any{"name": "Jane"}, map[string]any{"id": 1})
func Update[T any](table string, data map[string]any, where map[string]any) (*types.QueryResult[T], error) {
	return database.Update[T](table, data, where)
}

// Delete executes a DELETE query
// Usage: result, err := services.Delete[User]("users", map[string]any{"id": 1})
func Delete[T any](table string, where map[string]any) (*types.QueryResult[T], error) {
	return database.Delete[T](table, where)
}

// Raw executes a raw SQL query
// Usage: products, err := services.Raw[Product]("SELECT * FROM products WHERE price > ? ORDER BY price DESC", 100)
func Raw[T any](sql string, args ...any) (*types.QueryResult[T], error) {
	return database.Raw[T](sql, args...)
}

// RawSingle executes a raw SQL query and returns a single result
// Usage: user, err := services.RawSingle[User]("SELECT * FROM users WHERE email = ?", "john@example.com")
func RawSingle[T any](sql string, args ...any) (*T, error) {
	result, err := database.Raw[T](sql, args...)
	if err != nil {
		return nil, err
	}

	return result.Single, nil
}

// Transaction executes multiple operations in a single database transaction
// Usage:
//
//	err := services.Transaction(ctx, func(tx *pg.Tx) error {
//	    // Your database operations here
//	    return nil
//	})
func Transaction(ctx context.Context, operations ...func(*pg.Tx) error) error {
	return database.Transaction(ctx, operations...)
}

// Advanced query builders

// Query creates a new QueryParams builder for complex queries
// Usage:
//
//	query := services.Query().
//	    SetOperation("select").
//	    SetTable("users").
//	    AddWhere("age", 18).
//	    AddOrder("created_at DESC").
//	    SetLimit(10)
//	result, err := services.DB[User](query)
func Query() *types.QueryParams {
	return types.NewQuery()
}

// SelectBuilder creates a SELECT query builder
func SelectBuilder[T any](table string) *SelectQueryBuilder[T] {
	return &SelectQueryBuilder[T]{
		query: types.NewQuery().SetOperation("select").SetTable(table),
	}
}

// InsertBuilder creates an INSERT query builder
func InsertBuilder[T any](table string) *InsertQueryBuilder[T] {
	return &InsertQueryBuilder[T]{
		query: types.NewQuery().SetOperation("insert").SetTable(table),
	}
}

// UpdateBuilder creates an UPDATE query builder
func UpdateBuilder[T any](table string) *UpdateQueryBuilder[T] {
	return &UpdateQueryBuilder[T]{
		query: types.NewQuery().SetOperation("update").SetTable(table),
	}
}

// DeleteBuilder creates a DELETE query builder
func DeleteBuilder[T any](table string) *DeleteQueryBuilder[T] {
	return &DeleteQueryBuilder[T]{
		query: types.NewQuery().SetOperation("delete").SetTable(table),
	}
}

// Query builders for fluent interface

type SelectQueryBuilder[T any] struct {
	query *types.QueryParams
}

func (b *SelectQueryBuilder[T]) Select(columns []string) *SelectQueryBuilder[T] {
	b.query.SetSelect(columns)
	return b
}

func (b *SelectQueryBuilder[T]) Where(key string, value any) *SelectQueryBuilder[T] {
	b.query.AddWhere(key, value)
	return b
}

func (b *SelectQueryBuilder[T]) WhereRaw(whereClause string, args ...any) *SelectQueryBuilder[T] {
	b.query.SetWhereRaw(whereClause, args...)
	return b
}

func (b *SelectQueryBuilder[T]) Join(join string) *SelectQueryBuilder[T] {
	b.query.AddJoin(join)
	return b
}

func (b *SelectQueryBuilder[T]) Order(order string) *SelectQueryBuilder[T] {
	b.query.AddOrder(order)
	return b
}

func (b *SelectQueryBuilder[T]) Limit(limit int) *SelectQueryBuilder[T] {
	b.query.SetLimit(limit)
	return b
}

func (b *SelectQueryBuilder[T]) Offset(offset int) *SelectQueryBuilder[T] {
	b.query.SetOffset(offset)
	return b
}

func (b *SelectQueryBuilder[T]) Timeout(timeout time.Duration) *SelectQueryBuilder[T] {
	b.query.SetTimeout(timeout)
	return b
}

func (b *SelectQueryBuilder[T]) Context(ctx context.Context) *SelectQueryBuilder[T] {
	b.query.SetContext(ctx)
	return b
}

func (b *SelectQueryBuilder[T]) Execute() (*types.QueryResult[T], error) {
	return database.ExecuteQuery[T](b.query)
}

func (b *SelectQueryBuilder[T]) First() (*T, error) {
	b.query.SetLimit(1)
	result, err := database.ExecuteQuery[T](b.query)
	if err != nil {
		return nil, err
	}
	return result.Single, nil
}

type InsertQueryBuilder[T any] struct {
	query *types.QueryParams
}

func (b *InsertQueryBuilder[T]) Data(data map[string]any) *InsertQueryBuilder[T] {
	b.query.SetData(data)
	return b
}

func (b *InsertQueryBuilder[T]) Set(key string, value any) *InsertQueryBuilder[T] {
	b.query.AddData(key, value)
	return b
}

func (b *InsertQueryBuilder[T]) Returning(columns ...string) *InsertQueryBuilder[T] {
	b.query.SetReturning(columns...)
	return b
}

func (b *InsertQueryBuilder[T]) OnConflict(strategy string) *InsertQueryBuilder[T] {
	b.query.SetOnConflict(strategy)
	return b
}

func (b *InsertQueryBuilder[T]) Execute() (*types.QueryResult[T], error) {
	return database.ExecuteQuery[T](b.query)
}

type UpdateQueryBuilder[T any] struct {
	query *types.QueryParams
}

func (b *UpdateQueryBuilder[T]) Set(key string, value any) *UpdateQueryBuilder[T] {
	b.query.AddData(key, value)
	return b
}

func (b *UpdateQueryBuilder[T]) Data(data map[string]any) *UpdateQueryBuilder[T] {
	b.query.SetData(data)
	return b
}

func (b *UpdateQueryBuilder[T]) Where(key string, value any) *UpdateQueryBuilder[T] {
	b.query.AddWhere(key, value)
	return b
}

func (b *UpdateQueryBuilder[T]) WhereRaw(whereClause string, args ...any) *UpdateQueryBuilder[T] {
	b.query.SetWhereRaw(whereClause, args...)
	return b
}

func (b *UpdateQueryBuilder[T]) Returning(columns ...string) *UpdateQueryBuilder[T] {
	b.query.SetReturning(columns...)
	return b
}

func (b *UpdateQueryBuilder[T]) Execute() (*types.QueryResult[T], error) {
	return database.ExecuteQuery[T](b.query)
}

type DeleteQueryBuilder[T any] struct {
	query *types.QueryParams
}

func (b *DeleteQueryBuilder[T]) Where(key string, value any) *DeleteQueryBuilder[T] {
	b.query.AddWhere(key, value)
	return b
}

func (b *DeleteQueryBuilder[T]) WhereRaw(whereClause string, args ...any) *DeleteQueryBuilder[T] {
	b.query.SetWhereRaw(whereClause, args...)
	return b
}

func (b *DeleteQueryBuilder[T]) Returning(columns ...string) *DeleteQueryBuilder[T] {
	b.query.SetReturning(columns...)
	return b
}

func (b *DeleteQueryBuilder[T]) Execute() (*types.QueryResult[T], error) {
	return database.ExecuteQuery[T](b.query)
}
