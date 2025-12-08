package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sem1-final-project-hard-level/internal/config"
	custommiddleware "sem1-final-project-hard-level/internal/custom_middlewares"
	"sem1-final-project-hard-level/internal/dto"
	"sem1-final-project-hard-level/internal/handlers"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	priceHandler := handlers.NewPriceHandler()

	r.Route("/api", func(r chi.Router) {
		r.Route("/v0", func(r chi.Router) {
			r.With(custommiddleware.QueryParserMiddleware[dto.GetPricesQueryParamsDto](nil))
			r.Route("/prices", func(r chi.Router) { r.Get("/", priceHandler.GetPrices) })
		})
	})

	addr := fmt.Sprintf(":%d", a.config.Port)

	a.server = &http.Server{
		Addr:         addr,
		Handler:      r,
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
	return nil
}
