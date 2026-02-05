package handler

import (
	"net/http"

	"github.com/WardJune/with-chi/internal/transport"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	transport.Success(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
