package bootstrap

import (
	"log"

	"github.com/artumont/DotSlashStream/backend/internal/service/tmdb"
)

type Services struct {
	Tmdb *tmdb.Service
}

func (app *Application) SetupServices() {
	config := app.Config
	env := app.Env

	app.Services.Tmdb = tmdb.NewTmdbService(config.TMDBService.TMDBAPIUrl, config.TMDBService.TMDBAPIKey, env.ContextTimeout)
	log.Println("Services initialized successfully.")
}
