package server

import (
	"net/http"
	"time"

	"github.com/WardJune/with-chi/internal/config"
)

func New() *http.Server {

	cfg := config.Load()
	router := NewRouter()

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
