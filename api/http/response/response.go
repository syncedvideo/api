package response

import (
	"encoding/json"
	"net/http"
)

func WithJSON(w http.ResponseWriter, v interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

type errorResponse struct {
	Message string `json:"message"`
}

func WithError(w http.ResponseWriter, msg string, code int) error {
	err := errorResponse{
		Message: msg,
	}
	return WithJSON(w, err, code)
}
