package app

import (
	custommiddleware "sem1-final-project-hard-level/internal/custom_middlewares"
	"sem1-final-project-hard-level/internal/dto"
	"sem1-final-project-hard-level/internal/handlers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func GetChiRouter() *chi.Mux {
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
			r.Use(custommiddleware.QueryParserMiddleware[dto.GetPricesQueryParamsDto](nil))
			r.Get("/prices", priceHandler.GetPrices)
		})
	})

	return r
}
