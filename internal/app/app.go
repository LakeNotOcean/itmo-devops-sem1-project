package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sem1-final-project-hard-level/internal/config"
	"sem1-final-project-hard-level/internal/database"
)

type App struct {
	config *config.Config
	server *http.Server
}

func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a *App) Run() error {

	addr := fmt.Sprintf(":%d", a.config.Port)

	if err := database.InitDb(a.config); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := os.MkdirAll(a.config.TempFileDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory %s: %w", a.config.TempFileDir, err)
	}

	a.server = &http.Server{
		Addr:         addr,
		Handler:      GetChiRouter(a.config),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  time.Duration(a.config.IdleTimeout) * time.Second,
	}

	serverErr := make(chan error, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server starting on %s\n", addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("failed to start server: %w", err)
	case <-stop:
		log.Println("Shutting down server gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Println("Server stopped")
		return nil
	}
}

func (a *App) Close() error {
	log.Println("Cleaning up resources...")
	if err := database.CloseDb(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}
