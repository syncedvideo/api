package syncedvideo

import (
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
	id   string
	name string
}

const JSONContentType = "application/json"

func NewServer(store RoomStore) *Server {
	server := &Server{}
	server.store = store

	router := chi.NewMux()

	router.Get("/rooms/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", JSONContentType)
		id := chi.URLParam(r, "id")

		room := server.store.GetRoom(id)
		if room.id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(room.name))
	})

	server.Handler = router

	return server
}
