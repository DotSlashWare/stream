package bootstrap

import (
	"log"
	"os"
	"path/filepath"
)

// Initializes the application for the first time by setting up the databases and creating a flag file.
func (app *Application) FirstTimeSetup() {
	log.Println("Starting first-time setup.")

	app.Postgres.SetupProtocol()

	app.MarkFirstTimeSetupComplete()
	log.Println("First-time setup completed successfully.")
}

// Checks if the application is being set up for the first time by looking for a flag file.
func (app *Application) MarkFirstTimeSetupComplete() {
	flagDir := "/run/flags"

	// Create the flags directory if it doesn't exist
	if err := os.MkdirAll(flagDir, 0755); err != nil {
		log.Printf("Error creating flags directory: %v.", err)
		return
	}

	flagPath := filepath.Join(flagDir, "setup.flag")
	file, err := os.Create(flagPath)
	if err != nil {
		log.Printf("Error creating first-time setup flag file: %v.", err)
		return
	}
	defer file.Close()
	log.Println("First-time setup completed successfully, flag file created.")
}

// Checks if the application is being set up for the first time by looking for a flag file.
func (app *Application) IsFirstTimeSetup() bool {
	flagPath := filepath.Join("/run/flags", "setup.flag")
	file, err := os.OpenFile(flagPath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return true // File does not exist, indicating first-time setup
		}
		log.Printf("Error checking first-time setup: %v.", err)
		return true // Assume first-time setup if there's an error opening the file
	}
	defer file.Close()
	return false
}
