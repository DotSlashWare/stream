package postgres

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/artumont/DotSlashStream/backend/pkg/postgresparser"
)

// Setup function for the PostgreSQL database, creates tables and schemas as needed.
func (m *Manager) SetupProtocol() error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Database setup panicked: %v", r)
		}
	}()

	log.Println("Setting up PostgreSQL database")

	// Ensure migration tracking table exists
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	scriptsDir := "/app/scripts/postgres"

	// Validate scripts directory
	if err := m.validateScriptsDirectory(scriptsDir); err != nil {
		return fmt.Errorf("scripts directory validation failed: %w", err)
	}

	scripts, err := m.discoverScripts(scriptsDir)
	if err != nil {
		return fmt.Errorf("failed to discover scripts: %w", err)
	}

	if len(scripts) == 0 {
		log.Println("No setup scripts found in scripts directory")
		return nil
	}

	log.Printf("Found %d setup scripts", len(scripts))

	// Execute scripts in order
	for _, script := range scripts {
		if err := m.executeScript(script); err != nil {
			return fmt.Errorf("setup failed at script %s: %w", script.Name, err)
		}
	}

	log.Println("Database setup completed successfully.")
	return nil
}

// Creates the migration tracking table if it doesn't exist
func (m *Manager) ensureMigrationTable() error {
	migrationManager := postgresparser.NewMigrationManager(m.Client)
	return migrationManager.EnsureMigrationTable()
}

// Finds and parses all setup scripts in the directory
func (m *Manager) discoverScripts(scriptsDir string) ([]*postgresparser.ScriptInfo, error) {
	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read scripts directory: %w", err)
	}

	var scripts []*postgresparser.ScriptInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if !strings.HasSuffix(strings.ToLower(filename), ".sql") {
			continue
		}

		script, err := m.parseScript(scriptsDir, filename)
		if err != nil {
			log.Printf("Warning: failed to parse script %s: %v", filename, err)
			continue
		}

		scripts = append(scripts, script)
	}

	// Sort scripts by name to ensure consistent execution order
	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].Name < scripts[j].Name
	})

	return scripts, nil
}

// Reads and parses a script file
func (m *Manager) parseScript(scriptsDir, filename string) (*postgresparser.ScriptInfo, error) {
	filePath := filepath.Join(scriptsDir, filename)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	script := &postgresparser.ScriptInfo{
		Name:    filename,
		Path:    filePath,
		Content: string(content),
	}

	// Parse metadata from comments at the top of the file using postgresparser
	parser := postgresparser.NewParser()
	metadata := parser.ParseMetadata(string(content))
	if metadata != nil {
		metadata.Name = filename
		script.Metadata = metadata
		script.Dependencies = metadata.Dependencies
	}

	return script, nil
}

// Executes a single setup script
func (m *Manager) executeScript(script *postgresparser.ScriptInfo) error {
	migrationManager := postgresparser.NewMigrationManager(m.Client)

	// Check if script has already been executed successfully
	status, err := migrationManager.GetMigrationStatus(script.Name)
	if err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if status == "success" {
		log.Printf("Script %s already executed successfully, skipping", script.Name)
		return nil
	}

	log.Printf("Executing script: %s", script.Name)

	if script.Metadata != nil && script.Metadata.Description != "" {
		log.Printf("Description: %s", script.Metadata.Description)
	}

	// Check dependencies
	if err := migrationManager.CheckDependencies(script.Dependencies); err != nil {
		return fmt.Errorf("dependency check failed: %w", err)
	}

	// Execute the script in a transaction
	tx, err := m.Client.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Use the parser to execute the script
	parser := postgresparser.NewParser()
	result := parser.ExecuteScript(tx, script.Content)

	// Record execution result
	metadata := postgresparser.ScriptMetadata{
		Name:       script.Name,
		ExecutedAt: time.Now(),
	}

	if script.Metadata != nil {
		metadata.Description = script.Metadata.Description
		metadata.Version = script.Metadata.Version
		metadata.Author = script.Metadata.Author
		metadata.Dependencies = script.Metadata.Dependencies
	}

	if result.Success {
		metadata.Status = "success"
		log.Printf("Successfully executed script: %s", script.Name)
	} else {
		metadata.Status = "failed"
		metadata.Error = result.Error.Error()
		log.Printf("Failed to execute script %s: %v", script.Name, result.Error)
	}

	// Calculate checksum using SHA256 for reliable integrity checking
	hash := sha256.New()
	hash.Write([]byte(script.Content))
	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	// Record migration in the same transaction
	err = migrationManager.RecordMigration(tx, metadata, checksum)
	if err != nil {
		log.Printf("Warning: failed to record migration for %s: %v", script.Name, err)
	}

	if !result.Success {
		return result.Error
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Validates that the scripts directory exists and is a directory
func (m *Manager) validateScriptsDirectory(scriptsDir string) error {
	info, err := os.Stat(scriptsDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", scriptsDir)
	}
	if err != nil {
		return fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", scriptsDir)
	}
	return nil
}
