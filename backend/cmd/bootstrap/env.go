package bootstrap

import (
	"os"
	"strconv"
)

type Env struct {
	Version        string
	ContextTimeout int
	Port           string
	Environment    string
	PostgresURL    string
	ConfigPath     string
	FirstTimeSetup bool
}

func NewEnv() *Env {
	return &Env{
		Version:        getEnvOrDefault("VERSION", "1.0.0"),
		ContextTimeout: getEnvOrDefaultInt("CONTEXT_TIMEOUT", 30),
		Port:           getEnvOrDefault("PORT", "9340"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "release"),
		ConfigPath:     getEnvOrDefault("CONFIG_PATH", "/var/config/stream"),
		PostgresURL:    getEnvOrDefault("POSTGRES_URL", "postgresql://stream_admin:stream_admin@postgres:5432/stream?sslmode=disable"),
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
