package postgresparser

import (
	"time"
)

// ScriptMetadata represents metadata about a setup script
type ScriptMetadata struct {
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	Version      string    `json:"version,omitempty"`
	Author       string    `json:"author,omitempty"`
	Dependencies []string  `json:"dependencies,omitempty"`
	ExecutedAt   time.Time `json:"executed_at"`
	Status       string    `json:"status"`
	Error        string    `json:"error,omitempty"`
}

// ScriptInfo represents a discovered script
type ScriptInfo struct {
	Name         string
	Path         string
	Content      string
	Metadata     *ScriptMetadata
	Dependencies []string
}

// ScriptResult represents the result of script execution
type ScriptResult struct {
	Success bool
	Output  interface{}
	Error   error
}

// SQLStatement represents a parsed SQL statement
type SQLStatement struct {
	Type    string // CREATE, INSERT, UPDATE, DELETE, etc.
	Content string // The actual SQL
	LineNum int    // Line number in the script
}

// SQLScript represents a parsed SQL script with multiple statements
type SQLScript struct {
	Statements []SQLStatement
	Variables  map[string]string // For environment variable substitution
}
