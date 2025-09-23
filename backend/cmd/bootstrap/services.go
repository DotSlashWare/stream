package bootstrap

import (
	"log"
)

type Services struct {
}

func (app *Application) SetupServices() {
	// env := app.Env
	// config := app.Config

	log.Println("Services initialized successfully.")
}
