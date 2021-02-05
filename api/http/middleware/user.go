package middleware

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/request"
	"github.com/syncedvideo/syncedvideo/http/response"
)

const userCookieKey string = "userID"

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getUserFromCookie(r)
		if err != nil {
			response.WithError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		request.WithUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func getUserFromCookie(r *http.Request) (syncedvideo.User, error) {
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

	user, err := syncedvideo.Config.Store.User().Get(userID)
	if err != nil {
		return syncedvideo.User{}, err
	}
	return user, nil
}
