package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/artumont/DotSlashStream/backend/pkg/postgresparser"
)

// SelectFrom performs a SELECT query on the specified table with given columns and optional WHERE clause.
func (manager *Manager) SelectFrom(table string, columns []string, where string, args ...interface{}) (*sql.Rows, error) {
	if !validateString(table) || !validateStringList(columns) || (len(args) > 0 && !validateString(where)) {
		return nil, fmt.Errorf("invalid table name, columns, or where clause")
	}

	ctx, cancel := context.WithTimeout(context.Background(), manager.ContextTimeout)
	defer cancel()

	cols := strings.Join(columns, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", cols, table)
	if where != "" {
		query += " WHERE " + where
	}

	if err := postgresparser.ValidateSQL(query); err != nil {
		log.Printf("SQL validation error: %v", err)
		return nil, fmt.Errorf("sql validation: %w", err)
	}

	return manager.Client.QueryContext(ctx, query, args...)
}

// RunQuerySecure executes a SQL query securely after validating it (not meant for data return).
func (manager *Manager) RunQuerySecure(query string, args ...interface{}) (sql.Result, error) {
	if manager == nil || manager.Client == nil {
		return nil, fmt.Errorf("postgres manager or client is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), manager.ContextTimeout)
	defer cancel()

	if err := postgresparser.ValidateSQL(query); err != nil {
		log.Printf("SQL validation error: %v", err)
		return nil, fmt.Errorf("sql validation: %w", err)
	}

	res, err := manager.Client.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, fmt.Errorf("exec failed: %w", err)
	}

	return res, nil
}

// validateString checks if the input string adecuate for use in SQL queries
func validateString(input string) bool {
	inputContent := strings.TrimSpace(input)
	return inputContent != ""
}

// validateStringList checks if all strings in the input list are adecuate for use in SQL queries
func validateStringList(inputs []string) bool {
	if len(inputs) == 0 {
		return false
	}

	for _, item := range inputs {
		if !validateString(item) {
			return false
		}
	}

	return true
}
