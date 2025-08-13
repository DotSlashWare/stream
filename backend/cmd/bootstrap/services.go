package bootstrap

import (
	"log"

	"github.com/artumont/DotSlashStream/backend/internal/service/tmdb"
)

type Services struct {
	Tmdb *tmdb.Service
}

func (app *Application) SetupServices() {
	env := app.Env

	app.Services.Tmdb = tmdb.NewTmdbService(env.TheMovieDBAPI, env.TheMovieDBKey, env.ContextTimeout)
	log.Println("Services initialized successfully.")
}
