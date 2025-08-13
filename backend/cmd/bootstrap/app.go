package bootstrap

import (
	"log"
	"time"

	"github.com/artumont/DotSlashStream/backend/internal/database/postgres"
	"github.com/gin-gonic/gin"
)

// Encapsulates the entire application state, including environment settings, router, and database managers.
type Application struct {
	InitTime time.Time
	Env      *Env
	Router   *gin.Engine
	Postgres *postgres.Manager
	Services Services
}

// Shutdown protocol for gracefully shutting down the application. It closes all database connections and logs the shutdown process.
func (app *Application) Shutdown() {
	log.Println("Shutting down application...")
	// @logic: Postgres Shutdown
	if app.Postgres != nil {
		if err := app.Postgres.Close(); err != nil {
			log.Printf("Error closing Postgres connection: %v", err)
		}
	}
	log.Println("Application shutdown complete.")
}

// Factory function for creating a new Application instance.
func NewApplication() *Application {
	env := NewEnv()
	router := gin.New()
	err := router.SetTrustedProxies(nil) // Disable trusted proxies for security
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	return &Application{
		InitTime: time.Now(),
		Env:      env,
		Router:   router,
		Postgres: nil, // will be initialized in the init protocol
	}
}

func (app *Application) Start() {
	env := app.Env
	app.SetupLogging()
	app.SetupDatabases()
	app.RegisterMiddleware()
	app.RegisterControllers()

	if app.IsFirstTimeSetup() {
		app.FirstTimeSetup()
	} else {
		log.Println("Skipping first-time setup, application has been initialized before.")
	}

	app.Router.Run(env.GetServerAddress())
}
