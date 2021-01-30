package syncedvideo

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Handlers interface {
	User() UserHandler
	Room() RoomHandler
	UserMiddleware(http.Handler) http.Handler
	CorsMiddleware(http.Handler) http.Handler
}

type RoomHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Vote(w http.ResponseWriter, r *http.Request)
	Connect(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	Auth(w http.ResponseWriter, r *http.Request)
}

func RegisterHandlers(r chi.Router, h Handlers) {
	r.Use(h.CorsMiddleware)
	r.Route("/user", func(rr chi.Router) {
		rr.Post("/auth", h.User().Auth)
	})
	r.Route("/room", func(rr chi.Router) {
		rr.Use(h.UserMiddleware)
		rr.Post("/", h.Room().Create)
		rr.Get("/{roomID}", h.Room().Get)
		rr.Put("/{roomID}", h.Room().Update)
		rr.HandleFunc("/{roomID}/connect", h.Room().Connect)
	})
}
