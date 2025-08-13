package bootstrap

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"io"
	"log"
	"os"
	"time"
)

// Logging initialization protocol for setting up the app's logging with rotation.
func (app *Application) SetupLogging() {
	logPath := app.Env.LogPath

	writer, err := rotatelogs.New(
		logPath+"app.%Y-%m-%d-%H_%M",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	if err != nil {
		log.Fatalf("Failed to create log file writer: %v", err)
	}
	log.Println("Log file writer created successfully")

	multiWriter := io.MultiWriter(os.Stdout, writer)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Log output set to both stdout and file")

	log.Println("Logging setup complete")
}
