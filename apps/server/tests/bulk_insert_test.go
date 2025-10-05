package tests

import (
	"testing"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
)

func TestBulkInsertSQL(t *testing.T) {
	// Test data
	entries := []any{
		map[string]any{
			"id":       "user1",
			"username": "john",
			"email":    "john@example.com",
		},
		map[string]any{
			"id":       "user2",
			"username": "jane",
			"email":    "jane@example.com",
		},
	}

	columns := []string{"id", "username", "email"}

	// Test basic bulk insert SQL generation
	sql, values, err := database.BuildBulkInsertSQL("users", entries, columns, nil, "")
	if err != nil {
		t.Fatalf("buildBulkInsertSQL failed: %v", err)
	}

	expectedSQL := "INSERT INTO users (id, username, email) VALUES (?, ?, ?), (?, ?, ?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedValueCount := 6 // 2 entries * 3 columns each
	if len(values) != expectedValueCount {
		t.Errorf("Expected %d values, got %d", expectedValueCount, len(values))
	}

	// Validate the actual values
	expectedValues := []any{"user1", "john", "john@example.com", "user2", "jane", "jane@example.com"}
	for i, val := range values {
		if val != expectedValues[i] {
			t.Errorf("Expected value %v at index %d, got %v", expectedValues[i], i, val)
		}
	}

	// Test with RETURNING clause
	sql, _, err = database.BuildBulkInsertSQL("users", entries, columns, []string{"id", "username"}, "")
	if err != nil {
		t.Fatalf("buildBulkInsertSQL with RETURNING failed: %v", err)
	}

	expectedSQL = "INSERT INTO users (id, username, email) VALUES (?, ?, ?), (?, ?, ?) RETURNING id, username"
	if sql != expectedSQL {
		t.Errorf("Expected SQL with RETURNING: %s, got: %s", expectedSQL, sql)
	}

	// Test with ON CONFLICT
	sql, _, err = database.BuildBulkInsertSQL("users", entries, columns, nil, "(email) DO UPDATE SET username = EXCLUDED.username")
	if err != nil {
		t.Fatalf("buildBulkInsertSQL with ON CONFLICT failed: %v", err)
	}

	expectedSQL = "INSERT INTO users (id, username, email) VALUES (?, ?, ?), (?, ?, ?) ON CONFLICT (email) DO UPDATE SET username = EXCLUDED.username"
	if sql != expectedSQL {
		t.Errorf("Expected SQL with ON CONFLICT: %s, got: %s", expectedSQL, sql)
	}

	// Test error cases
	_, _, err = database.BuildBulkInsertSQL("", entries, columns, nil, "")
	if err == nil {
		t.Error("Expected error for empty table name")
	}

	_, _, err = database.BuildBulkInsertSQL("users", []any{}, columns, nil, "")
	if err == nil {
		t.Error("Expected error for empty entries")
	}

	_, _, err = database.BuildBulkInsertSQL("users", entries, []string{}, nil, "")
	if err == nil {
		t.Error("Expected error for empty columns")
	}

	// Test invalid entry type
	invalidEntries := []any{"not a map"}
	_, _, err = database.BuildBulkInsertSQL("users", invalidEntries, columns, nil, "")
	if err == nil {
		t.Error("Expected error for invalid entry type")
	}
}

func TestExtractColumnsFromEntries(t *testing.T) {
	entries := []any{
		map[string]any{
			"id":       "user1",
			"username": "john",
			"email":    "john@example.com",
		},
		map[string]any{
			"id":       "user2",
			"username": "jane",
			"role":     "admin", // Additional column in second entry
		},
	}

	columns, err := database.ExtractColumnsFromEntries(entries)
	if err != nil {
		t.Fatalf("extractColumnsFromEntries failed: %v", err)
	}

	expectedColumns := map[string]bool{
		"id":       true,
		"username": true,
		"email":    true,
		"role":     true,
	}

	if len(columns) != len(expectedColumns) {
		t.Errorf("Expected %d columns, got %d", len(expectedColumns), len(columns))
	}

	for _, col := range columns {
		if !expectedColumns[col] {
			t.Errorf("Unexpected column: %s", col)
		}
	}
}

func TestQueryValidation(t *testing.T) {
	// Test valid bulk insert
	query := types.NewQuery().
		SetOperation("insert").
		SetTable(lib.TableUsers).
		SetEntries([]any{
			map[string]any{"id": "1", "name": "test"},
		})

	err := query.Validate()
	if err != nil {
		t.Errorf("Valid bulk insert query should not fail validation: %v", err)
	}

	// Test valid single insert
	query2 := types.NewQuery().
		SetOperation("insert").
		SetTable(lib.TableUsers).
		SetData(map[string]any{"id": "1", "name": "test"})

	err = query2.Validate()
	if err != nil {
		t.Errorf("Valid single insert query should not fail validation: %v", err)
	}

	// Test invalid query with both Data and Entries
	query3 := types.NewQuery().
		SetOperation("insert").
		SetTable(lib.TableUsers).
		SetData(map[string]any{"id": "1", "name": "test"}).
		SetEntries([]any{map[string]any{"id": "2", "name": "test2"}})

	err = query3.Validate()
	if err == nil {
		t.Error("Query with both Data and Entries should fail validation")
	}

	// Test invalid query with neither Data nor Entries
	query4 := types.NewQuery().
		SetOperation("insert").
		SetTable(lib.TableUsers)

	err = query4.Validate()
	if err == nil {
		t.Error("Query with neither Data nor Entries should fail validation")
	}
}
