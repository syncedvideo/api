package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
)

// Handler
type Handler struct {
	*chi.Mux
	store syncedvideo.Store
	redis *redis.Client
}

// RegisterHandlers registers all handlers
func RegisterHandlers(store syncedvideo.Store, redis *redis.Client) *Handler {
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
	room := syncedvideo.Room{ID: uuid.New()}
	if err := h.store.Room().Create(&room); err != nil {
		log.Printf("error creating room: %s", err)
		http.Error(w, "error creating room", http.StatusInternalServerError)
		return
	}
	log.Printf("created room id: %v", room)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

func (h Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		log.Printf("error parsing uuid: %v", err)
		http.Error(w, "room id is invalid", 400)
		return
	}
	room, err := h.store.Room().Get(id)
	if err != nil {
		log.Printf("error getting room: %v", err)
		http.Error(w, "error getting room", 400)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
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
