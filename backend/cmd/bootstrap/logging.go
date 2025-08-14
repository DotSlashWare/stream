package bootstrap

import (
	"log"
	"os"
)

// Logging initialization protocol for setting up the app's logging.
func (app *Application) SetupLogging() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("[APP] ")
	log.Println("Logging setup complete.")
}
