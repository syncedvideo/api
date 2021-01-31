package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/response"
)

const userCookieKey string = "userID"

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var userContextKey contextKey = contextKey("user")

func UserMiddleware(next http.Handler, userStore syncedvideo.UserStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromCookie(r, userStore)
		if err != nil {
			response.WithError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	})
}

func GetUserCtx(r *http.Request) syncedvideo.User {
	return r.Context().Value(userContextKey).(syncedvideo.User)
}

func getUserFromCookie(r *http.Request, userStore syncedvideo.UserStore) (syncedvideo.User, error) {
	userIDCookie, err := r.Cookie(userCookieKey)
	if err != nil {
		return syncedvideo.User{}, err
	}
	if userIDCookie.Value == "" {
		return syncedvideo.User{}, errors.New(userCookieKey + " cookie value is empty")
	}

	userID, err := uuid.Parse(userIDCookie.Value)
	if err != nil {
		return syncedvideo.User{}, err
	}
	if userID == uuid.Nil {
		return syncedvideo.User{}, errors.New("userID is nil")
	}

	user, err := userStore.Get(userID)
	if err != nil {
		return syncedvideo.User{}, err
	}
	return user, nil
}
