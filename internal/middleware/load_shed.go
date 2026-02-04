package middleware

import (
	"net/http"
	"time"

	"github.com/WardJune/with-chi/pkg/metrics"
)

// LoadShedding creates a middleware that limits the number of concurrent requests to a maximum.
// If the maximum is reached, it waits for a specified duration before allowing a new request.
func LoadShedding(max int, wait time.Duration) func(http.Handler) http.Handler {
	sem := make(chan struct{}, max)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case sem <- struct{}{}:
				metrics.InFlight.Inc()
				defer func() {
					<-sem
					metrics.InFlight.Dec()
				}()

				start := time.Now()
				next.ServeHTTP(w, r)
				metrics.RequestDuration.Observe(float64(time.Since(start).Seconds()))
			case <-r.Context().Done():
				return
			case <-time.After(wait):
				metrics.ShedTotal.Inc()
				w.Header().Set("Retry-After", "5")
				http.Error(w, "server overload", http.StatusServiceUnavailable)
			}
		})
	}
}
