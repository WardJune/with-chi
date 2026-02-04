package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WardJune/with-chi/internal/config"
	"github.com/WardJune/with-chi/internal/server"
	"github.com/WardJune/with-chi/pkg/metrics"
)

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Info("Server forced to shutdown with error:", "error", err.Error())
	}

	slog.Info("Server exiting")

	done <- true
}

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	srv := server.New(cfg)
	metrics.Register()

	done := make(chan bool, 1)

	go gracefulShutdown(srv, done)

	slog.Info("Server started", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
	}

	<-done
	slog.Info("shutting down server")
}
