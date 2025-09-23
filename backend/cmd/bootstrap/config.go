package bootstrap

import (
	"log"
	"path/filepath"

	"github.com/artumont/DotSlashStream/backend/internal/utils/config"
)

func (app *Application) LoadConfig() {
	env := app.Env
	configPath := filepath.Join(env.ConfigPath, "backend.config.json")

	var err error
	app.Config, err = config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	log.Printf("Config loaded successfully from %s.", configPath)
}