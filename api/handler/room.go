package handler

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/syncedvideo/syncedvideo"
)

// RoomHandler implements the synvedvideo.RoomHandler interface
type RoomHandler struct {
	store syncedvideo.Store
	redis *redis.Client
}

func NewRoomHandler(s syncedvideo.Store, r *redis.Client) syncedvideo.RoomHandler {
	return &RoomHandler{
		store: s,
		redis: r,
	}
}
func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *RoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *RoomHandler) Connect(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *RoomHandler) Vote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}
