package server

import (
	"github.com/WardJune/with-chi/internal/handler"
	"github.com/WardJune/with-chi/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)

	r.Get("/health", handler.HealthHandler)

	return r
}
