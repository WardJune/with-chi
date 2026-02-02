package handler

import (
	"net/http"
	"time"

	"github.com/WardJune/with-chi/internal/transport"
)

type HelloWorldResponse struct {
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(500 * time.Millisecond)

	transport.Success(w, 200, HelloWorldResponse{
		Message: "Hello World",
		Status:  "success",
	})
}
