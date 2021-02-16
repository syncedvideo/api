package middleware

import (
	"net/http"
)

const (
	allowOriginHeaderKey      = "Access-Control-Allow-Origin"
	allowMethodsHeaderKey     = "Access-Control-Allow-Methods"
	allowHeadersHeaderKey     = "Access-Control-Allow-Headers"
	allowCredentialsHeaderKey = "Access-Control-Allow-Credentials"

	allowedMethods = "GET,HEAD,OPTIONS,POST,PUT"
	allowedHeaders = "Content-Type"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(allowOriginHeaderKey, r.Header.Get("Origin"))
		w.Header().Set(allowMethodsHeaderKey, allowedMethods)
		w.Header().Set(allowHeadersHeaderKey, allowedHeaders)
		w.Header().Set(allowCredentialsHeaderKey, "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
