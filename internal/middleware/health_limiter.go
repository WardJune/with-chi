package middleware

import (
	"net/http"
	"time"

	"github.com/WardJune/with-chi/internal/limiter"
)

func HealthLimiter(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.TryAcquire(500 * time.Millisecond) {
				http.Error(w, "busy", http.StatusServiceUnavailable)
				return
			}
			defer l.Release()
			next.ServeHTTP(w, r)
		})
	}
}
