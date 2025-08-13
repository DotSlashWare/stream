package bootstrap

import (
	"log"

	"github.com/artumont/DotSlashStream/backend/internal/middleware/logger"
)

// Controller registration protocol for setting up route controllers.
func (app *Application) RegisterControllers() {
	// env := app.Env

	log.Println("Registered controllers successfully")
}

// Middleware registration protocol for setting up global middleware.
func (app *Application) RegisterMiddleware() {
	env := app.Env

	{ // @logic: Logger Middleware (global)
		logger := logger.NewLoggerMiddleware(app.Router, env.LogPath)
		logger.Register(app.Router)
		log.Printf("Logger middleware registered with log path: %s", env.LogPath)
	}

	log.Println("Registered global middleware successfully")
}
