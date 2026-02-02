package handler

import (
	"net/http"

	"github.com/WardJune/with-chi/internal/transport"
)

type HelloWorldResponse struct {
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	transport.Success(w, 200, HelloWorldResponse{
		Message: "Hello World",
		Status:  "success",
	})
}
