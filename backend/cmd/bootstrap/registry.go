package bootstrap

import (
	"log"

	"github.com/artumont/DotSlashStream/backend/internal/controller/health"
	"github.com/artumont/DotSlashStream/backend/internal/middleware/logger"
)

// Controller registration protocol for setting up route controllers.
func (app *Application) RegisterControllers() {
	// env := app.Env
	// config := app.Config

	baseGroup := app.Router.Group("/")
	{
		healthController := health.NewHealthController(app.InitTime, app.Postgres)
		healthController.Register(baseGroup)
	}

	apiGroup := app.Router.Group("/api")
	{
		movieController := movie.NewMovieController(app.Services.Tmdb)
		movieController.Register(apiGroup)
	}

	log.Println("Registered controllers successfully.")
}

// Middleware registration protocol for setting up global middleware.
func (app *Application) RegisterMiddleware() {
	// env := app.Env

	{ // @logic: Logger Middleware (global)
		logger := logger.NewLoggerMiddleware(app.Router)
		logger.Register(app.Router)
		log.Println("Logger middleware registered")
	}

	log.Println("Registered global middleware successfully.")
}
