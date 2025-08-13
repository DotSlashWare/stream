package postgresparser

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Handles migration tracking and rollback operations
type MigrationManager struct {
	db *sql.DB
}

// Creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// Creates the migration tracking table if it doesn't exist
func (m *MigrationManager) EnsureMigrationTable() error {
	createTableSQL := `
		CREATE SCHEMA IF NOT EXISTS migrations;

		CREATE TABLE IF NOT EXISTS migrations.script_migrations (
		name VARCHAR(255) PRIMARY KEY,
		description TEXT,
		version VARCHAR(50),
		author VARCHAR(100),
		dependencies TEXT[], -- Array of dependency script names
		executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failed', 'skipped')),
		error_message TEXT,
		rollback_sql TEXT, -- Optional rollback SQL for the migration
		checksum VARCHAR(64), -- SHA256 checksum of the script content
		CONSTRAINT valid_status CHECK (status IN ('success', 'failed', 'skipped'))
		);

		CREATE INDEX IF NOT EXISTS idx_script_migrations_executed_at ON migrations.script_migrations(executed_at);
		CREATE INDEX IF NOT EXISTS idx_script_migrations_status ON migrations.script_migrations(status);
		CREATE INDEX IF NOT EXISTS idx_script_migrations_checksum ON migrations.script_migrations(checksum);
	`

	_, err := m.db.Exec(createTableSQL)
	return err
}

// Records a migration execution in the tracking table
func (m *MigrationManager) RecordMigration(tx *sql.Tx, metadata ScriptMetadata, checksum string) error {
	// Convert dependencies to PostgreSQL array format
	var dependenciesArray interface{}
	if len(metadata.Dependencies) > 0 {
		// Convert to PostgreSQL array format
		dependenciesArray = fmt.Sprintf("{%s}", strings.Join(metadata.Dependencies, ","))
	} else {
		dependenciesArray = "{}"
	}

	_, err := tx.Exec(`
		INSERT INTO migrations.script_migrations 
		(name, description, version, author, dependencies, executed_at, status, error_message, checksum)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (name) DO UPDATE SET
		description = EXCLUDED.description,
		version = EXCLUDED.version,
		author = EXCLUDED.author,
		dependencies = EXCLUDED.dependencies,
		executed_at = EXCLUDED.executed_at,
		status = EXCLUDED.status,
		error_message = EXCLUDED.error_message,
		checksum = EXCLUDED.checksum
		`,
		metadata.Name,
		metadata.Description,
		metadata.Version,
		metadata.Author,
		dependenciesArray,
		metadata.ExecutedAt,
		metadata.Status,
		metadata.Error,
		checksum,
	)

	return err
}

// Checks if a migration has been executed successfully
func (m *MigrationManager) GetMigrationStatus(name string) (string, error) {
	var status string
	err := m.db.QueryRow(
		"SELECT status FROM migrations.script_migrations WHERE name = $1",
		name,
	).Scan(&status)

	if err == sql.ErrNoRows {
		return "", nil // Migration not found
	}

	return status, err
}

// Verifies that all required dependencies have been executed
func (m *MigrationManager) CheckDependencies(dependencies []string) error {
	if len(dependencies) == 0 {
		return nil
	}

	for _, dep := range dependencies {
		var status string
		err := m.db.QueryRow(
			"SELECT status FROM migrations.script_migrations WHERE name = $1 AND status = 'success'",
			dep,
		).Scan(&status)

		if err != nil {
			return fmt.Errorf("dependency %s not found or not successfully executed", dep)
		}
	}

	return nil
}

// Returns the migration history
func (m *MigrationManager) GetMigrationHistory() ([]ScriptMetadata, error) {
	rows, err := m.db.Query(`
		SELECT name, description, version, author, executed_at, status, error_message
		FROM migrations.script_migrations 
		ORDER BY executed_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []ScriptMetadata
	for rows.Next() {
		var migration ScriptMetadata
		var description, version, author, errorMsg sql.NullString

		err := rows.Scan(
			&migration.Name,
			&description,
			&version,
			&author,
			&migration.ExecutedAt,
			&migration.Status,
			&errorMsg,
		)
		if err != nil {
			return nil, err
		}

		// Handle nullable fields
		if description.Valid {
			migration.Description = description.String
		}
		if version.Valid {
			migration.Version = version.String
		}
		if author.Valid {
			migration.Author = author.String
		}
		if errorMsg.Valid {
			migration.Error = errorMsg.String
		}

		migrations = append(migrations, migration)
	}

	return migrations, rows.Err()
}

// Performs a rollback of a specific migration
func (m *MigrationManager) RollbackMigration(name string, rollbackSQL string) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin rollback transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute rollback SQL if provided
	if rollbackSQL != "" {
		parser := NewParser()
		result := parser.ExecuteScript(tx, rollbackSQL)
		if !result.Success {
			return fmt.Errorf("rollback SQL failed: %w", result.Error)
		}
	}

	// Update migration status
	_, err = tx.Exec(`
		UPDATE migrations.script_migrations 
		SET status = 'rolled_back', executed_at = $1
		WHERE name = $2
	`, time.Now(), name)

	if err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}

	return tx.Commit()
}

// Returns migrations that need to be executed
func (m *MigrationManager) GetPendingMigrations(availableScripts []string) ([]string, error) {
	if len(availableScripts) == 0 {
		return nil, nil
	}

	// Create a placeholder string for the IN clause
	placeholders := make([]string, len(availableScripts))
	args := make([]interface{}, len(availableScripts))
	for i, script := range availableScripts {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = script
	}

	query := fmt.Sprintf(`
		SELECT name FROM migrations.script_migrations 
		WHERE name IN (%s) AND status = 'success'
	`, strings.Join(placeholders, ","))

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect executed migrations
	executedMigrations := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		executedMigrations[name] = true
	}

	// Find pending migrations
	var pending []string
	for _, script := range availableScripts {
		if !executedMigrations[script] {
			pending = append(pending, script)
		}
	}

	return pending, rows.Err()
}

// Checks if migrations have been tampered with
func (m *MigrationManager) ValidateMigrationIntegrity(scripts map[string]string) error {
	// Get all executed migrations with checksums
	rows, err := m.db.Query(`
		SELECT name, checksum FROM migrations.script_migrations 
		WHERE status = 'success' AND checksum IS NOT NULL
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, dbChecksum string
		if err := rows.Scan(&name, &dbChecksum); err != nil {
			return err
		}

		// Check if we have the script content
		if content, exists := scripts[name]; exists {
			currentChecksum := calculateChecksum(content)
			if currentChecksum != dbChecksum {
				return fmt.Errorf("migration %s has been modified (checksum mismatch)", name)
			}
		}
	}

	return rows.Err()
}
