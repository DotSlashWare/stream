package bootstrap

import (
	"os"
	"strconv"
)

type Env struct {
	ContextTimeout int
	Port           string
	LogPath        string
	Environment    string
	PostgresURL    string
	TheMovieDBAPI  string
	TheMovieDBKey  string
	FirstTimeSetup bool
}

func NewEnv() *Env {
	return &Env{
		ContextTimeout: getEnvOrDefaultInt("CONTEXT_TIMEOUT", 30),
		Port:           getEnvOrDefault("PORT", "9340"),
		LogPath:        getEnvOrDefault("LOG_PATH", "/var/log/stream"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "release"),
		PostgresURL:    getEnvOrDefault("POSTGRES_URL", "postgresql://stream_admin:stream_admin@postgres:5432/stream?sslmode=disable"),
		TheMovieDBAPI:  getEnvOrDefault("TMDB_API", "https://api.themoviedb.org/3"),
		TheMovieDBKey:  getEnvOrDefault("TMDB_KEY", "none"),
	}
}

func (env *Env) GetServerAddress() string {
	return "0.0.0.0" + ":" + env.Port
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}

	return defaultValue
}

/*
func getEnvOrDefaultBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.ParseBool(value); err == nil {
			return v
		}
	}

	return defaultValue
}
*/
