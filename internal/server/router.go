package server

import (
	"time"

	"github.com/WardJune/with-chi/internal/handler"
	"github.com/WardJune/with-chi/internal/limiter"
	"github.com/WardJune/with-chi/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	rLimiter := limiter.NewAdaptiveLimiter(10, 300*time.Millisecond, 0.2)

	healthLimiter := limiter.NewLimiter(2)

	//Sub-router
	metricRouter := chi.NewRouter()
	metricRouter.Get("/", promhttp.Handler().ServeHTTP)

	healthRouter := chi.NewRouter()
	healthRouter.Use(middleware.HealthLimiter(healthLimiter))
	healthRouter.Get("/", handler.HealthHandler)

	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AdaptiveShedding(rLimiter))
		r.Get("/", handler.HelloWorldHandler)
	})

	r.Mount("/metrics", metricRouter)
	r.Mount("/health", healthRouter)

	return r
}
