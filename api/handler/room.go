package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/syncedvideo/syncedvideo"
)

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
	user := GetUserCtx(r)
	room := syncedvideo.Room{OwnerUserID: user.ID}
	if err := h.store.Room().Create(&room); err != nil {
		log.Printf("error creating room: %s\n", err)
		RespondWithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	log.Printf("created room id: %v\n", room.ID)
	RespondWithJSON(w, room, http.StatusCreated)
}

func (h *RoomHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		log.Printf("error parsing uuid: %v", err)
		RespondWithError(w, "room not found", http.StatusNotFound)
		return
	}
	room, err := h.store.Room().Get(id)
	if err == sql.ErrNoRows {
		RespondWithError(w, "room not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("error getting room: %v", err)
		RespondWithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	RespondWithJSON(w, room, http.StatusOK)
}

func (h *RoomHandler) Update(w http.ResponseWriter, r *http.Request) {

	panic("not implemented") // TODO: Implement
}

func (h *RoomHandler) Vote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *RoomHandler) Connect(w http.ResponseWriter, r *http.Request) {
	// get room
	roomID, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		log.Printf("error parsing roomID: %v\n", err)
		return
	}
	if roomID == uuid.Nil {
		log.Panicln("roomID is nil")
		return
	}
	room, err := h.store.Room().Get(roomID)
	if err != nil {
		log.Printf("error getting room: %v\n", err)
		return
	}
	log.Printf("connect to room id: %v\n", room.ID)

	// upgrade http to tcp
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading to websocket: %v\n", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			break
		}
		log.Printf("recieved message: %v\n", msg)
	}
}
