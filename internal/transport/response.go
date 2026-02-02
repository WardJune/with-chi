package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func DecodeJson(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func Success(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}

	fmt.Println(data)
}

func Error(w http.ResponseWriter, status int, code, message string) {
	Success(w, status, map[string]string{
		"code":    code,
		"message": message,
	})
}

func ValidationError(w http.ResponseWriter, err error) {
	fields := map[string]string{}

	for _, e := range err.(validator.ValidationErrors) {
		fields[e.Field()] = e.Tag()
	}

	Success(w, http.StatusBadRequest, map[string]any{
		"code":   "VALIDATION_ERROR",
		"fields": fields,
	})
}
