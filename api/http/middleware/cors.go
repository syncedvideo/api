package middleware

import (
	"net/http"

	"github.com/syncedvideo/syncedvideo/http/response"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			response.WithError(w, "missing Origin header", http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}
