package bootstrap

import (
	"log"

	"github.com/artumont/DotSlashStream/backend/internal/controller/health"
	"github.com/artumont/DotSlashStream/backend/internal/controller/movie"
	"github.com/artumont/DotSlashStream/backend/internal/middleware/logger"
)

// Controller registration protocol for setting up route controllers.
func (app *Application) RegisterControllers() {
	// env := app.Env
	// config := app.Config

	{ // @logic: Health Controller

		healthController := health.NewHealthController(app.InitTime, app.Postgres)
		healthController.Register(app.Router)
	}

	{ // @logic: Movie Controller
		movieController := movie.NewMovieController(app.Services.Tmdb)
		movieController.Register(app.Router)
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
