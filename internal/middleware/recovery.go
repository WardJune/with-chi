package middleware

import (
	"net/http"

	"github.com/WardJune/with-chi/internal/transport"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				transport.Error(w, http.StatusInternalServerError, "PANIC", "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
