package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/artumont/DotSlashStream/backend/internal/utils/sanitizer"
	_ "github.com/lib/pq"
)

// Struct for managing PostgreSQL connections. It holds the client, context timeout, and connection time.
type Manager struct {
	Client         *sql.DB
	ContextTimeout time.Duration
	ConnectionTime time.Time
	LogSanitizer   *sanitizer.Sanitizer
}

// Factory function for creating a new Manager instance. It connects to PostgreSQL using the provided connection string and context timeout, and performs a ping to ensure the connection is established.
func NewManager(connectionString string, contextTimout int, logSanitizer *sanitizer.Sanitizer) *Manager {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(contextTimout)*time.Second)
	defer cancel()

	log.Printf("Connecting to PostgreSQL at %s with timeout %d seconds.", logSanitizer.CleanString(connectionString), contextTimout)
	client, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v.", err)
	}
	log.Println("Connected to PostgreSQL successfully.")

	log.Printf("Pinging PostgreSQL at %s.", logSanitizer.CleanString(connectionString))
	if err = client.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v.", err)
	}
	log.Println("PostgreSQL ping successful.")

	return &Manager{
		Client:         client,
		ContextTimeout: time.Duration(contextTimout) * time.Second,
		ConnectionTime: time.Now(),
		LogSanitizer:   logSanitizer,
	}
}

// Performs a health check on the PostgreSQL connection by pinging the server.
func (manager *Manager) IsHealthy() bool {
	ctx, cancel := context.WithTimeout(context.Background(), manager.ContextTimeout)
	defer cancel()

	err := manager.Client.PingContext(ctx)
	if err != nil {
		log.Printf("PostgreSQL is not healthy: %v.", err)
		return false
	}

	log.Println("PostgreSQL is healthy.")
	return true
}

// Closes the PostgreSQL connection gracefully, ensuring all resources are released.
func (manager *Manager) Close() error {
	err := manager.Client.Close()
	if err != nil {
		log.Printf("Failed to close PostgreSQL connection: %v.", err)
		return err
	}

	log.Println("PostgreSQL connection closed successfully.")
	return nil
}
