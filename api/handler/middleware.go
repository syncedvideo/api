package handler

import (
	"context"
	"net/http"

	"github.com/syncedvideo/syncedvideo"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var userContextKey contextKey = contextKey("user")

func (h *Handlers) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromCookie(r, h.store.User())
		if err != nil {
			RespondWithError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	})
}

func GetUserCtx(r *http.Request) syncedvideo.User {
	return r.Context().Value(userContextKey).(syncedvideo.User)
}

func (h *Handlers) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
