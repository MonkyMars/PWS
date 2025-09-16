package types

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// QueryParams represents the parameters for building dynamic database queries.
// It supports all CRUD operations (Create, Read, Update, Delete) plus raw SQL execution.
// This struct provides a flexible and reusable way to interact with the database.
type QueryParams struct {
	// Operation specifies the type of database operation to perform
	// Supported values: "select", "insert", "update", "delete", "raw"
	Operation string `json:"operation"`

	// Table specifies the target table name (optional for raw SQL)
	Table string `json:"table,omitempty"`

	// Data contains the data to insert or update
	Data map[string]any `json:"data,omitempty"`

	// Select specifies which columns to select (for SELECT operations)
	Select []string `json:"select,omitempty"`

	// Where contains the WHERE clause conditions
	Where map[string]any `json:"where,omitempty"`

	// WhereRaw allows for complex WHERE conditions with raw SQL
	WhereRaw string `json:"where_raw,omitempty"`

	// WhereArgs contains arguments for WhereRaw placeholders
	WhereArgs []any `json:"where_args,omitempty"`

	// Join contains JOIN clauses
	Join []string `json:"join,omitempty"`

	// Order contains ORDER BY clauses
	Order []string `json:"order,omitempty"`

	// GroupBy contains GROUP BY clauses
	GroupBy []string `json:"group_by,omitempty"`

	// Having contains HAVING clause conditions
	Having map[string]any `json:"having,omitempty"`

	// Limit specifies the maximum number of records to return
	Limit int `json:"limit,omitempty"`

	// Offset specifies the number of records to skip
	Offset int `json:"offset,omitempty"`

	// RawSQL contains raw SQL query (for raw operations)
	RawSQL string `json:"raw_sql,omitempty"`

	// RawArgs contains arguments for RawSQL placeholders
	RawArgs []any `json:"raw_args,omitempty"`

	// Context for the database operation (optional)
	Context context.Context `json:"-"`

	// Timeout for the operation (optional)
	Timeout time.Duration `json:"timeout,omitempty"`

	// Transaction specifies if this operation should run in a transaction
	UseTransaction bool `json:"use_transaction,omitempty"`

	// Returning specifies columns to return (for INSERT/UPDATE/DELETE with RETURNING)
	Returning []string `json:"returning,omitempty"`

	// OnConflict specifies conflict resolution for INSERT operations
	OnConflict string `json:"on_conflict,omitempty"`

	// Additional parameters for custom use cases
	Params map[string]any `json:"params,omitempty"`
}

// QueryResult represents the result of a database operation
type QueryResult[T any] struct {
	// Data contains the result data
	Data []T `json:"data,omitempty"`

	// Single contains a single result (for operations returning one record)
	Single *T `json:"single,omitempty"`

	// Count contains the number of affected/returned rows
	Count int64 `json:"count"`

	// Success indicates if the operation was successful
	Success bool `json:"success"`

	// Error contains any error that occurred
	Error error `json:"error,omitempty"`

	// ExecutionTime contains the time taken to execute the query
	ExecutionTime time.Duration `json:"execution_time"`

	// Query contains the actual SQL query that was executed (for debugging)
	Query string `json:"query,omitempty"`

	// Args contains the arguments used in the query (for debugging)
	Args []any `json:"args,omitempty"`
}

// NewQuery creates a new QueryParams with default values
func NewQuery() *QueryParams {
	return &QueryParams{
		Where:  make(map[string]any),
		Data:   make(map[string]any),
		Having: make(map[string]any),
		Params: make(map[string]any),
	}
}

// SetOperation sets the operation type
func (q *QueryParams) SetOperation(operation string) *QueryParams {
	q.Operation = operation
	return q
}

// SetTable sets the target table
func (q *QueryParams) SetTable(table string) *QueryParams {
	q.Table = table
	return q
}

// SetData sets the data for insert/update operations
func (q *QueryParams) SetData(data map[string]any) *QueryParams {
	q.Data = data
	return q
}

// AddData adds a single key-value pair to the data
func (q *QueryParams) AddData(key string, value any) *QueryParams {
	if q.Data == nil {
		q.Data = make(map[string]any)
	}
	q.Data[key] = value
	return q
}

// SetSelect sets the columns to select
func (q *QueryParams) SetSelect(columns ...string) *QueryParams {
	q.Select = columns
	return q
}

// AddWhere adds a WHERE condition
func (q *QueryParams) AddWhere(key string, value any) *QueryParams {
	if q.Where == nil {
		q.Where = make(map[string]any)
	}
	q.Where[key] = value
	return q
}

// SetWhereRaw sets a raw WHERE clause
func (q *QueryParams) SetWhereRaw(whereClause string, args ...any) *QueryParams {
	q.WhereRaw = whereClause
	q.WhereArgs = args
	return q
}

// AddJoin adds a JOIN clause
func (q *QueryParams) AddJoin(join string) *QueryParams {
	q.Join = append(q.Join, join)
	return q
}

// AddOrder adds an ORDER BY clause
func (q *QueryParams) AddOrder(order string) *QueryParams {
	q.Order = append(q.Order, order)
	return q
}

// SetLimit sets the LIMIT
func (q *QueryParams) SetLimit(limit int) *QueryParams {
	q.Limit = limit
	return q
}

// SetOffset sets the OFFSET
func (q *QueryParams) SetOffset(offset int) *QueryParams {
	q.Offset = offset
	return q
}

// SetRawSQL sets raw SQL query
func (q *QueryParams) SetRawSQL(sql string, args ...any) *QueryParams {
	q.RawSQL = sql
	q.RawArgs = args
	q.Operation = "raw"
	return q
}

// SetTimeout sets operation timeout
func (q *QueryParams) SetTimeout(timeout time.Duration) *QueryParams {
	q.Timeout = timeout
	return q
}

// SetContext sets the context
func (q *QueryParams) SetContext(ctx context.Context) *QueryParams {
	q.Context = ctx
	return q
}

// WithTransaction enables transaction for this operation
func (q *QueryParams) WithTransaction() *QueryParams {
	q.UseTransaction = true
	return q
}

// SetReturning sets columns to return
func (q *QueryParams) SetReturning(columns ...string) *QueryParams {
	q.Returning = columns
	return q
}

// SetOnConflict sets conflict resolution
func (q *QueryParams) SetOnConflict(strategy string) *QueryParams {
	q.OnConflict = strategy
	return q
}

// Validate checks if the QueryParams is valid for the specified operation
func (q *QueryParams) Validate() error {
	switch strings.ToLower(q.Operation) {
	case "select":
		// Select operations don't require additional validation
		return nil
	case "insert":
		if len(q.Data) == 0 {
			return ErrNoDataProvided
		}
		return nil
	case "update":
		if len(q.Data) == 0 {
			return ErrNoDataProvided
		}
		if len(q.Where) == 0 && q.WhereRaw == "" {
			return ErrNoWhereCondition
		}
		return nil
	case "delete":
		if len(q.Where) == 0 && q.WhereRaw == "" {
			return ErrNoWhereCondition
		}
		return nil
	case "raw":
		if q.RawSQL == "" {
			return ErrNoRawSQL
		}
		return nil
	default:
		return ErrInvalidOperation
	}
}

// Common query builder errors
var (
	ErrNoDataProvided   = fmt.Errorf("no data provided for insert/update operation")
	ErrNoWhereCondition = fmt.Errorf("no WHERE condition provided for update/delete operation")
	ErrNoRawSQL         = fmt.Errorf("no raw SQL provided for raw operation")
	ErrInvalidOperation = fmt.Errorf("invalid operation specified")
)
