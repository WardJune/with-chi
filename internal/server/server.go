package server

import (
	"net/http"
	"time"

	"github.com/WardJune/with-chi/internal/config"
)

func New(cfg config.Config) *http.Server {
	router := NewRouter()

	return &http.Server{
		Addr:         "0.0.0.0:" + cfg.Port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
