package main

import (
	"log"
	"os"

	"sem1-final-project-hard-level/internal/app"
	"sem1-final-project-hard-level/internal/config"
)

func main() {
	cfg := config.Load()
	app := app.New(cfg)

	if err := app.Run(); err != nil {
		log.Printf("Fatal error: %v\n", err)

		if closeErr := app.Close(); closeErr != nil {
			log.Printf("Error during cleanup: %v", closeErr)
		}
		os.Exit(1)
	}

	if err := app.Close(); err != nil {
		log.Printf("Error during cleanup: %v\n", err)
		os.Exit(1)
	}
}
