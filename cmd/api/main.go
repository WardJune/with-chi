package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/WardJune/with-chi/internal/config"
	"github.com/WardJune/with-chi/internal/limiter"
	"github.com/WardJune/with-chi/internal/server"
	"github.com/WardJune/with-chi/pkg/metrics"
)

var (
	connNew    atomic.Int64
	connActive atomic.Int64
	connIdle   atomic.Int64
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

	go func() {
		for range time.Tick(2 * time.Second) {
			slog.Info("runtime", "goroutine", runtime.NumGoroutine())
			slog.Info("conn", "new", connNew.Load(), "active", connActive.Load(), "idle", connIdle.Load())
		}
	}()

	srv := server.New(cfg)
	extraMetrics := limiter.Metrics()
	metrics.Register(extraMetrics...)

	// srv.ConnState = func(conn net.Conn, state http.ConnState) {
	// 	switch state {
	// 	case http.StateNew:
	// 		connNew.Add(1)
	// 	case http.StateActive:
	// 		connActive.Add(1)
	// 	case http.StateIdle:
	// 		connIdle.Add(1)
	// 	}
	// }
	// srv.SetKeepAlivesEnabled(false)

	done := make(chan bool, 1)

	go gracefulShutdown(srv, done)

	slog.Info("Server started", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
	}

	<-done
	slog.Info("shutting down server")
}
