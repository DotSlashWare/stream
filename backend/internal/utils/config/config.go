package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	TMDBService      TMDBServiceConfig      `json:"tmdb_service"`
	InvidiousService InvidiousServiceConfig `json:"invidious_service"`
	LocalService     LocalServiceConfig     `json:"local_service"`
}

func NewConfig() *Config {
	return &Config{
		TMDBService: TMDBServiceConfig{
			TMDBAPIUrl: "https://api.themoviedb.org/3",
			TMDBAPIKey: "",
		},
		InvidiousService: InvidiousServiceConfig{
			VideoAPIUrl: "https://invidious.example.com/api/v1",
			VideoAPIKey: "",
		},
		LocalService: LocalServiceConfig{
			MediaPath: "/var/media/stream",
		},
	}
}

func LoadConfig(configPath string) (*Config, error) {
	var configData *Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, creating new config.", configPath)
		configData = NewConfig()
		configData.SetupProtocol()
		return configData, nil
	}

	fileData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v.", err)
		return nil, err
	}

	if err := json.Unmarshal(fileData, &configData); err != nil {
		log.Fatalf("Failed to unmarshal config JSON: %v.", err)
		return nil, err
	}

	return configData, nil
}
