package bootstrap

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/artumont/DotSlashStream/backend/internal/utils/config"
)

func (app *Application) LoadConfig() {
	log.Println("Loading configuration.")
    configPath := filepath.Join("/var/config/stream", "backend.config.json")

    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        log.Printf("Config file not found at %s, creating new config.", configPath)
        app.Config = config.NewConfig()
        app.Config.SetupProtocol()
        return
    }

    fileData, err := os.ReadFile(configPath)
    if err != nil {
        log.Fatalf("Failed to read config file: %v.", err)
        
    }

    var configData config.Config
    if err := json.Unmarshal(fileData, &configData); err != nil {
        log.Fatalf("Failed to unmarshal config JSON: %v.", err)
    }

    app.Config = &configData
    log.Printf("Config loaded successfully from %s.", configPath)
}