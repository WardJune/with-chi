package middleware

import (
	"net/http"
	"time"

	"github.com/WardJune/with-chi/internal/limiter"
	"github.com/WardJune/with-chi/pkg/metrics"
)

func AdaptiveShedding(l *limiter.AdaptiveLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.TryAcquire() {
				metrics.ShedTotal.Inc()
				http.Error(w, "overload", http.StatusServiceUnavailable)
			}

			metrics.InFlight.Inc()

			start := time.Now()
			defer func() {
				metrics.InFlight.Dec()
				metrics.RequestDuration.Observe(float64(time.Since(start).Seconds()))
				latency := time.Since(start)
				l.Obeserve(latency)
				l.Adjust()
				l.Release()
			}()

			next.ServeHTTP(w, r)
		})
	}
}
