package bootstrap

import (
	"log"

	"github.com/artumont/DotSlashStream/backend/internal/database/postgres"
	"github.com/artumont/DotSlashStream/backend/internal/utils/sanitizer"
	"github.com/artumont/DotSlashStream/backend/internal/utils/sanitizer/filters"
)

// Database initialization protocol for setting up the Postgres connection.
func (app *Application) SetupDatabases() {
	env := app.Env

	logSanitizer := sanitizer.NewSanitizer()
	logSanitizer.AddFilter(filters.NewDatabaseFilter())

	app.Postgres = postgres.NewManager(env.PostgresURL, env.ContextTimeout, logSanitizer)
	log.Println("Databases initialized successfully.")
}
