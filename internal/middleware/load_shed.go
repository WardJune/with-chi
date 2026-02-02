package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LoadShedding creates a middleware that limits the number of concurrent requests to a maximum.
// If the maximum is reached, it waits for a specified duration before allowing a new request.
func LoadShedding(max int, wait time.Duration) func(http.Handler) http.Handler {
	sem := make(chan struct{}, max)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
				next.ServeHTTP(w, r)
			case <-r.Context().Done():
				return
			case <-time.After(wait):
				w.Header().Set("Retry-After", "5")
				fmt.Println("server overload")
				http.Error(w, "server overload", http.StatusServiceUnavailable)
			}
		})
	}
}
