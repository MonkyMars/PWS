package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MonkyMars/PWS/types"
	"github.com/go-pg/pg/v10"
)

// DB executes database operations based on QueryParams and returns typed results.
// This function supports all CRUD operations (Create, Read, Update, Delete) plus raw SQL.
//
// Usage examples:
//
//	result, err := Database[User](query.NewQuery().SetOperation("select").AddWhere("id", 1))
//	result, err := Database[User](query.NewQuery().SetOperation("insert").AddData("name", "John"))
//	result, err := Database[Product](query.NewQuery().SetRawSQL("SELECT * FROM products WHERE price > ?", 100))
func Database[T any](query *types.QueryParams) (*types.QueryResult[T], error) {
	start := time.Now()
	result := &types.QueryResult[T]{
		Success: false,
	}

	// Validate the query
	if err := query.Validate(); err != nil {
		result.Error = err
		result.ExecutionTime = time.Since(start)
		return result, err
	}

	// Get database instance
	db := GetInstance()
	if db == nil {
		err := fmt.Errorf("database instance not initialized")
		result.Error = err
		result.ExecutionTime = time.Since(start)
		return result, err
	}

	// Set up context
	ctx := query.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// Apply timeout if specified
	if query.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, query.Timeout)
		defer cancel()
	}

	// Execute operation based on type
	var err error
	switch strings.ToLower(query.Operation) {
	case "select":
		err = executeSelect(ctx, db, query, result)
	case "insert":
		err = executeInsert(ctx, db, query, result)
	case "update":
		err = executeUpdate(ctx, db, query, result)
	case "delete":
		err = executeDelete(ctx, db, query, result)
	case "raw":
		err = executeRaw(ctx, db, query, result)
	default:
		err = fmt.Errorf("unsupported operation: %s", query.Operation)
	}

	// Set final result properties
	result.ExecutionTime = time.Since(start)
	result.Success = err == nil
	if err != nil {
		result.Error = err
	}

	return result, err
}

// executeSelect handles SELECT operations
func executeSelect[T any](ctx context.Context, db *DB, query *types.QueryParams, result *types.QueryResult[T]) error {
	var data []T
	var single T

	// Build the query
	pgQuery := db.ModelContext(ctx, &data)

	// Apply table if specified
	if query.Table != "" {
		pgQuery = pgQuery.Table(query.Table)
	}

	// Apply SELECT columns
	if len(query.Select) > 0 {
		for _, col := range query.Select {
			pgQuery = pgQuery.Column(col)
		}
	}

	// Apply WHERE conditions
	pgQuery = applyWhereConditions(pgQuery, query)

	// Apply JOINs
	for _, join := range query.Join {
		pgQuery = pgQuery.Join(join)
	}

	// Apply GROUP BY
	for _, groupBy := range query.GroupBy {
		pgQuery = pgQuery.Group(groupBy)
	}

	// Apply HAVING
	for key, value := range query.Having {
		pgQuery = pgQuery.Having("? = ?", pg.Ident(key), value)
	}

	// Apply ORDER BY
	for _, order := range query.Order {
		pgQuery = pgQuery.Order(order)
	}

	// Apply LIMIT
	if query.Limit > 0 {
		pgQuery = pgQuery.Limit(query.Limit)
	}

	// Apply OFFSET
	if query.Offset > 0 {
		pgQuery = pgQuery.Offset(query.Offset)
	}

	// Execute query
	err := pgQuery.Select()
	if err != nil {
		return fmt.Errorf("failed to execute select query: %w", err)
	}

	result.Data = data
	result.Count = int64(len(data))

	// Set single result if only one record
	if len(data) == 1 {
		result.Single = &data[0]
	} else if len(data) == 0 && query.Limit == 1 {
		// Try to get single record directly
		singleQuery := db.ModelContext(ctx, &single)
		if query.Table != "" {
			singleQuery = singleQuery.Table(query.Table)
		}

		// Apply same conditions
		singleQuery = applyWhereConditions(singleQuery, query)

		err = singleQuery.First()
		if err != nil && err != pg.ErrNoRows {
			return fmt.Errorf("failed to execute single select query: %w", err)
		}

		if err != pg.ErrNoRows {
			result.Single = &single
			result.Data = []T{single}
			result.Count = 1
		}
	}

	return nil
}

// executeInsert handles INSERT operations
func executeInsert[T any](ctx context.Context, db *DB, query *types.QueryParams, result *types.QueryResult[T]) error {
	// Build a raw INSERT query since go-pg's Value method works differently
	if query.Table == "" {
		return fmt.Errorf("table name is required for insert operation")
	}

	// Build column and value lists
	var columns []string
	var values []any
	for key, value := range query.Data {
		columns = append(columns, key)
		values = append(values, value)
	}

	if len(columns) == 0 {
		return fmt.Errorf("no data provided for insert")
	}

	// Build the SQL
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		query.Table,
		joinStrings(columns, ", "),
		buildPlaceholders(len(values)))

	// Add RETURNING clause if specified
	if len(query.Returning) > 0 {
		sql += " RETURNING " + joinStrings(query.Returning, ", ")
	}

	// Add ON CONFLICT clause if specified
	if query.OnConflict != "" {
		sql += " ON CONFLICT " + query.OnConflict
	}

	// Execute the query
	if len(query.Returning) > 0 {
		// If RETURNING is specified, we expect data back
		var returnedData []T
		_, err := db.QueryContext(ctx, &returnedData, sql, values...)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
		result.Count = int64(len(returnedData))
		result.Data = returnedData
		if len(returnedData) > 0 {
			result.Single = &returnedData[0]
		}
	} else {
		// No RETURNING clause
		res, err := db.ExecContext(ctx, sql, values...)
		if err != nil {
			return fmt.Errorf("failed to execute insert query: %w", err)
		}
		result.Count = int64(res.RowsAffected())
	}

	return nil
}

// executeUpdate handles UPDATE operations
func executeUpdate[T any](ctx context.Context, db *DB, query *types.QueryParams, result *types.QueryResult[T]) error {
	var data T

	// Build the query
	pgQuery := db.ModelContext(ctx, &data)

	// Apply table if specified
	if query.Table != "" {
		pgQuery = pgQuery.Table(query.Table)
	}

	// Apply WHERE conditions
	pgQuery = applyWhereConditions(pgQuery, query)

	// Apply data updates
	for key, value := range query.Data {
		pgQuery = pgQuery.Set("? = ?", pg.Ident(key), value)
	}

	// Apply RETURNING columns
	if len(query.Returning) > 0 {
		for _, col := range query.Returning {
			pgQuery = pgQuery.Returning(col)
		}
	}

	// Execute update
	res, err := pgQuery.Update()
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	result.Count = int64(res.RowsAffected())

	return nil
}

// executeDelete handles DELETE operations
func executeDelete[T any](ctx context.Context, db *DB, query *types.QueryParams, result *types.QueryResult[T]) error {
	var data T

	// Build the query
	pgQuery := db.ModelContext(ctx, &data)

	// Apply table if specified
	if query.Table != "" {
		pgQuery = pgQuery.Table(query.Table)
	}

	// Apply WHERE conditions
	pgQuery = applyWhereConditions(pgQuery, query)

	// Apply RETURNING columns
	if len(query.Returning) > 0 {
		for _, col := range query.Returning {
			pgQuery = pgQuery.Returning(col)
		}
	}

	// Execute delete
	res, err := pgQuery.Delete()
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	result.Count = int64(res.RowsAffected())

	return nil
}

// executeRaw handles raw SQL operations
func executeRaw[T any](ctx context.Context, db *DB, query *types.QueryParams, result *types.QueryResult[T]) error {
	var data []T

	// Store the actual query for debugging
	result.Query = query.RawSQL
	result.Args = query.RawArgs

	// Execute raw query
	_, err := db.QueryContext(ctx, &data, query.RawSQL, query.RawArgs...)
	if err != nil {
		return fmt.Errorf("failed to execute raw query: %w", err)
	}

	result.Data = data
	result.Count = int64(len(data))

	// Set single result if only one record
	if len(data) == 1 {
		result.Single = &data[0]
	}

	return nil
}

// applyWhereConditions applies WHERE conditions to a query
func applyWhereConditions(pgQuery *pg.Query, query *types.QueryParams) *pg.Query {
	// Apply simple WHERE conditions
	for key, value := range query.Where {
		pgQuery = pgQuery.Where("? = ?", pg.Ident(key), value)
	}

	// Apply raw WHERE conditions
	if query.WhereRaw != "" {
		pgQuery = pgQuery.Where(query.WhereRaw, query.WhereArgs...)
	}

	return pgQuery
}

// Transaction executes multiple operations in a single transaction
func Transaction(ctx context.Context, operations ...func(*pg.Tx) error) error {
	db := GetInstance()
	if db == nil {
		return fmt.Errorf("database instance not initialized")
	}

	return db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		for _, operation := range operations {
			if err := operation(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

// Quick helper functions for common operations

// Select executes a SELECT query
func Select[T any](table string, where map[string]any) (*types.QueryResult[T], error) {
	query := types.NewQuery().
		SetOperation("select").
		SetTable(table)

	for key, value := range where {
		query.AddWhere(key, value)
	}

	return Database[T](query)
}

// Insert executes an INSERT query
func Insert[T any](table string, data map[string]any) (*types.QueryResult[T], error) {
	query := types.NewQuery().
		SetOperation("insert").
		SetTable(table).
		SetData(data)

	return Database[T](query)
}

// Update executes an UPDATE query
func Update[T any](table string, data map[string]any, where map[string]any) (*types.QueryResult[T], error) {
	query := types.NewQuery().
		SetOperation("update").
		SetTable(table).
		SetData(data)

	for key, value := range where {
		query.AddWhere(key, value)
	}

	return Database[T](query)
}

// Delete executes a DELETE query
func Delete[T any](table string, where map[string]any) (*types.QueryResult[T], error) {
	query := types.NewQuery().
		SetOperation("delete").
		SetTable(table)

	for key, value := range where {
		query.AddWhere(key, value)
	}

	return Database[T](query)
}

// Raw executes a raw SQL query
func Raw[T any](sql string, args ...any) (*types.QueryResult[T], error) {
	query := types.NewQuery().SetRawSQL(sql, args...)
	return Database[T](query)
}

// Helper functions
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func buildPlaceholders(count int) string {
	if count == 0 {
		return ""
	}
	if count == 1 {
		return "?"
	}

	result := "?"
	for i := 1; i < count; i++ {
		result += ", ?"
	}
	return result
}
