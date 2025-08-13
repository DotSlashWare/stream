package logger

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// Custom middleware for logging HTTP requests in a web application.
type LoggerMiddleware struct {
	logPath   string
	Formatter *CustomFormatter // Might need to be accessed in the future so it'll remain public for now
}

// Factory function to create a new LoggerMiddleware instance.
func NewLoggerMiddleware(router *gin.Engine, logPath string) *LoggerMiddleware {
	return &LoggerMiddleware{
		logPath:   logPath,
		Formatter: &CustomFormatter{},
	}
}

// Sets up the logging middleware for the Gin router with rotation.
func (m *LoggerMiddleware) Register(router *gin.Engine) {
	logPath := m.logPath

	writer, err := rotatelogs.New(
		logPath+"gin.%Y-%m-%d-%H_%M",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	if err != nil {
		log.Fatalf("Failed to create gin log file writer: %v", err)
	}
	log.Println("Gin Log file writer created successfully")

	multiWriter := io.MultiWriter(os.Stdout, writer)
	router.Use(gin.LoggerWithConfig(
		gin.LoggerConfig{
			Formatter: m.Formatter.Format,
			Output:    multiWriter,
			SkipPaths: []string{"/health", "/metrics"}, // Skip logging for health and metrics endpoints
		},
	))
	log.Println("Gin Logger middleware registered successfully")
}
