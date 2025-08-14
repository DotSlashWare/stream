package logger

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// Custom middleware for logging HTTP requests in a web application.
type LoggerMiddleware struct {
	Formatter *CustomFormatter // Might need to be accessed in the future so it'll remain public for now
}

// Factory function to create a new LoggerMiddleware instance.
func NewLoggerMiddleware(router *gin.Engine) *LoggerMiddleware {
	return &LoggerMiddleware{
		Formatter: &CustomFormatter{},
	}
}

// Sets up the logging middleware for the Gin router with rotation.
func (middleware *LoggerMiddleware) Register(router *gin.Engine) {
	router.Use(gin.LoggerWithConfig(
		gin.LoggerConfig{
			Formatter: middleware.Formatter.Format,
			Output:    os.Stdout,
			SkipPaths: []string{"/health", "/metrics"}, // Skip logging for health and metrics endpoints
		},
	))
	log.Println("Gin Logger middleware registered successfully for stdout output.")
}
