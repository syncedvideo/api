package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/syncedvideo/syncedvideo"
)

type Handlers struct {
	store syncedvideo.Store
	user  syncedvideo.UserHandler
	room  syncedvideo.RoomHandler
}

func New(s syncedvideo.Store, r *redis.Client) syncedvideo.Handlers {
	return &Handlers{
		store: s,
		user:  NewUserHandler(s),
		room:  NewRoomHandler(s, r),
	}
}

func (h *Handlers) User() syncedvideo.UserHandler {
	return h.user
}

func (h *Handlers) Room() syncedvideo.RoomHandler {
	return h.room
}

func RespondWithJSON(w http.ResponseWriter, v interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

type errorResponse struct {
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, msg string, code int) error {
	err := errorResponse{
		Message: msg,
	}
	return RespondWithJSON(w, err, code)
}
