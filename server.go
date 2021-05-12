package syncedvideo

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RoomStore interface {
	GetRoom(id string) Room
}

type Server struct {
	store RoomStore
	http.Handler
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

const jsonContentType = "application/json"

func NewServer(store RoomStore) *Server {
	server := new(Server)
	server.store = store

	router := chi.NewMux()

	router.Get("/rooms/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", jsonContentType)
		id := chi.URLParam(r, "id")

		room := server.store.GetRoom(id)
		if room.ID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(room)
	})

	server.Handler = router

	return server
}
