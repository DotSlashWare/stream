package postgresparser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Parser handles parsing and execution of PostgreSQL scripts
type Parser struct{}

// NewParser creates a new PostgreSQL script parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseMetadata extracts metadata from script comments
func (p *Parser) ParseMetadata(content string) *ScriptMetadata {
	lines := strings.Split(content, "\n")
	var metadataLines []string

	// Look for JSON metadata in comments at the start of the file
	inMetadata := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "-- METADATA:") {
			inMetadata = true
			continue
		}

		if inMetadata {
			if strings.HasPrefix(line, "--") {
				// Remove comment prefix and add to metadata
				metadataLine := strings.TrimPrefix(line, "--")
				metadataLine = strings.TrimSpace(metadataLine)
				if metadataLine != "" {
					metadataLines = append(metadataLines, metadataLine)
				}
			} else {
				// End of metadata section
				break
			}
		}
	}

	if len(metadataLines) == 0 {
		return nil
	}

	// Try to parse as JSON
	jsonStr := strings.Join(metadataLines, "")
	var metadata ScriptMetadata
	if err := json.Unmarshal([]byte(jsonStr), &metadata); err != nil {
		log.Printf("Warning: failed to parse script metadata: %v.", err)
		return nil
	}

	return &metadata
}

// ExecuteScript executes SQL content within a transaction
func (p *Parser) ExecuteScript(tx *sql.Tx, sqlContent string) ScriptResult {
	if len(strings.TrimSpace(sqlContent)) == 0 {
		return ScriptResult{
			Success: true,
			Output:  "Script is empty, skipped",
		}
	}

	// Parse the SQL script into statements
	script, err := p.ParseSQL(sqlContent)
	if err != nil {
		return ScriptResult{
			Success: false,
			Error:   fmt.Errorf("failed to parse SQL script: %w", err),
		}
	}

	// Execute each statement
	var results []interface{}
	for i, stmt := range script.Statements {
		result, err := p.executeStatement(tx, stmt)
		if err != nil {
			return ScriptResult{
				Success: false,
				Error:   fmt.Errorf("failed to execute statement %d (%s): %w", i+1, stmt.Type, err),
			}
		}
		results = append(results, result)
	}

	return ScriptResult{
		Success: true,
		Output:  results,
	}
}

// ParseSQL parses SQL content into individual statements
func (p *Parser) ParseSQL(sqlContent string) (*SQLScript, error) {
	script := &SQLScript{
		Variables: make(map[string]string),
	}

	// Clean the content
	content := p.cleanSQL(sqlContent)

	// Split into statements (basic approach - could be enhanced)
	statements := p.splitStatements(content)

	for i, stmt := range statements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		stmtType := p.detectStatementType(stmt)
		script.Statements = append(script.Statements, SQLStatement{
			Type:    stmtType,
			Content: strings.TrimSpace(stmt),
			LineNum: i + 1,
		})
	}

	return script, nil
}

// cleanSQL removes comments and normalizes whitespace
func (p *Parser) cleanSQL(content string) string {
	lines := strings.Split(content, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip comment lines that aren't metadata
		if strings.HasPrefix(line, "--") && !strings.HasPrefix(line, "-- METADATA:") {
			continue
		}
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

// splitStatements splits SQL content into individual statements
func (p *Parser) splitStatements(content string) []string {
	// Simple statement splitting by semicolon
	var statements []string
	var current strings.Builder
	inQuotes := false
	var quoteChar rune

	for _, char := range content {
		switch char {
		case '\'', '"':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
			}
			current.WriteRune(char)
		case ';':
			if !inQuotes {
				statements = append(statements, current.String())
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	// Add the last statement if it doesn't end with semicolon
	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}

// detectStatementType determines the type of SQL statement
func (p *Parser) detectStatementType(stmt string) string {
	stmt = strings.TrimSpace(strings.ToUpper(stmt))

	patterns := map[string]string{
		"^CREATE\\s+(TABLE|INDEX|SCHEMA|DATABASE|VIEW|FUNCTION|PROCEDURE|TRIGGER)": "CREATE",
		"^DROP\\s+(TABLE|INDEX|SCHEMA|DATABASE|VIEW|FUNCTION|PROCEDURE|TRIGGER)":   "DROP",
		"^ALTER\\s+(TABLE|INDEX|SCHEMA|DATABASE|VIEW)":                             "ALTER",
		"^INSERT\\s+INTO": "INSERT",
		"^UPDATE\\s+":     "UPDATE",
		"^DELETE\\s+FROM": "DELETE",
		"^SELECT\\s+":     "SELECT",
		"^GRANT\\s+":      "GRANT",
		"^REVOKE\\s+":     "REVOKE",
		"^SET\\s+":        "SET",
		"^COMMENT\\s+ON":  "COMMENT",
	}

	for pattern, stmtType := range patterns {
		if matched, _ := regexp.MatchString(pattern, stmt); matched {
			return stmtType
		}
	}

	return "UNKNOWN"
}

// executeStatement executes a single SQL statement
func (p *Parser) executeStatement(tx *sql.Tx, stmt SQLStatement) (interface{}, error) {
	switch stmt.Type {
	case "SELECT":
		return p.executeSelect(tx, stmt.Content)
	case "INSERT", "UPDATE", "DELETE":
		return p.executeModification(tx, stmt.Content)
	case "CREATE", "DROP", "ALTER", "GRANT", "REVOKE", "SET", "COMMENT":
		return p.executeDDL(tx, stmt.Content)
	default:
		return p.executeDDL(tx, stmt.Content) // Default to DDL execution
	}
}

// executeSelect executes a SELECT statement and returns results
func (p *Parser) executeSelect(tx *sql.Tx, sql string) (interface{}, error) {
	rows, err := tx.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}

// executeModification executes INSERT, UPDATE, DELETE statements
func (p *Parser) executeModification(tx *sql.Tx, sql string) (interface{}, error) {
	result, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"rows_affected": rowsAffected,
	}, nil
}

// executeDDL executes DDL statements (CREATE, DROP, ALTER, etc.)
func (p *Parser) executeDDL(tx *sql.Tx, sql string) (interface{}, error) {
	_, err := tx.Exec(sql)
	if err != nil {
		return nil, err
	}

	return "DDL executed successfully", nil
}

// ValidateSQL performs basic SQL validation
func (p *Parser) ValidateSQL(sqlContent string) error {
	// Basic validation - check for obvious issues
	content := strings.TrimSpace(sqlContent)
	if content == "" {
		return fmt.Errorf("SQL content is empty")
	}

	// Check for basic dangerous SQL patterns (very basic)
	dangerousPatterns := []string{
		"DROP DATABASE",
		"TRUNCATE",
		"DELETE FROM.*WHERE.*1.*=.*1",
	}

	upperContent := strings.ToUpper(content)
	for _, pattern := range dangerousPatterns {
		if matched, _ := regexp.MatchString(pattern, upperContent); matched {
			log.Printf("Warning: Potentially dangerous SQL pattern detected: %s.", pattern)
		}
	}

	return nil
}
