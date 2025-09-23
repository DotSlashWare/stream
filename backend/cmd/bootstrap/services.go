package bootstrap

import (
	"log"
)

type Services struct {
	TMDBService interface{}
}

func (app *Application) SetupServices() {
	// env := app.Env
	// config := app.Config

	log.Println("Services initialized successfully.")
}
