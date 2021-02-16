package request

import (
	"context"
	"net/http"

	"github.com/syncedvideo/syncedvideo"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	userContextKey contextKey = contextKey("user")
	roomContextKey contextKey = contextKey("room")
)

func WithUser(r *http.Request, user syncedvideo.User) {
	*r = *r.WithContext(context.WithValue(r.Context(), userContextKey, user))
}

func GetUserCtx(r *http.Request) syncedvideo.User {
	return r.Context().Value(userContextKey).(syncedvideo.User)
}

func WithRoom(r *http.Request, room syncedvideo.Room) {
	*r = *r.WithContext(context.WithValue(r.Context(), roomContextKey, room))
}

func GetRoomCtx(r *http.Request) syncedvideo.Room {
	return r.Context().Value(roomContextKey).(syncedvideo.Room)
}
