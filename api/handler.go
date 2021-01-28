package syncedvideo

import (
	"net/http"

	"github.com/go-chi/chi"
)

type RoomHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Vote(w http.ResponseWriter, r *http.Request)
	Connect(w http.ResponseWriter, r *http.Request)
}

func RegisterRoomHandler(r chi.Router, h RoomHandler) {
	r.Route("/room", func(rr chi.Router) {
		rr.Post("/", h.Create)
		rr.Get("/{roomID}", h.Get)
		rr.Put("/{roomID}", h.Update)
		rr.HandleFunc("/{roomID}/connect", h.Connect)
	})
}

type UserHandler interface {
	Auth(w http.ResponseWriter, r *http.Request)
}

func RegisterUserHandler(r chi.Router, h UserHandler) {
	r.Route("/user", func(rr chi.Router) {
		rr.Post("/auth", h.Auth)
	})
}
