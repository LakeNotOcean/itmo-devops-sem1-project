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

	defer func() {
		if closeErr := app.Close(); closeErr != nil {
			log.Printf("Error during cleanup: %v", closeErr)
		}
	}()

	if err := app.Run(); err != nil {
		log.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}
