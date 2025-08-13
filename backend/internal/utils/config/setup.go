package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func (c *Config) SetupProtocol() {
	configDir := "/var/config/stream"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Failed to create config directory: %v.", err)
		return
	}

	configPath := filepath.Join(configDir, "backend.config.json")

	configData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal config to JSON: %v.", err)
		return
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		log.Printf("Failed to write config file: %v.", err)
		return
	}

	log.Printf("Config file created successfully at %s.", configPath)
	log.Println("Server config setup complete.")
}