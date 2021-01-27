package syncedvideo

import (
	"net/http"

	"github.com/go-chi/chi"
)

type RoomHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Connect(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Vote(w http.ResponseWriter, r *http.Request)
}

func RegisterRoomHandler(m *chi.Mux, h RoomHandler) {
	m.Route("/room", func(r chi.Router) {
		m.Post("/", h.Create)
		m.Get("/{roomID}", h.Get)
		m.Put("/{roomID}", h.Update)
	})
}
