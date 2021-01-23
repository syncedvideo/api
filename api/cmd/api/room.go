package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/syncedvideo/syncedvideo"
)

// Handler
type Handler struct {
	*chi.Mux
	store api.Store
	redis *redis.Client
}

// RegisterHandlers registers all handlers
func RegisterHandlers(store api.Store, redis *redis.Client) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
		redis: redis,
	}
	h.Route("/room", func(r chi.Router) {
		r.Post("/", h.CreateRoom)
		r.Route("/{roomID}", func(r chi.Router) {
			r.Get("/", h.GetRoom)
			r.Put("/", h.UpdateRoom)
			r.Delete("/", h.DeleteRoom)
			r.Post("/player/resume", h.ResumePlayer)
			r.Post("/player/pause", h.PausePlayer)
			r.Post("/player/fast-forward", h.FastForwardPlayer)
			r.Post("/player/skip", h.SkipPlayer)
			r.Post("/queue/items", h.AddQueueItem)
			r.Delete("/queue/items/{queueItemID}", h.RemoveQueueItem)
			r.Post("/queue/items/{queueItemID}/vote", h.VoteQueueItem)
		})
	})
	return h
}

func (h Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) ResumePlayer(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) PausePlayer(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) FastForwardPlayer(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) SkipPlayer(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) AddQueueItem(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) RemoveQueueItem(w http.ResponseWriter, r *http.Request) {
	//
}

func (h Handler) VoteQueueItem(w http.ResponseWriter, r *http.Request) {
	//
}
