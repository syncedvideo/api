package handler

import (
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
