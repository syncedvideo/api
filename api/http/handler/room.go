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
	"github.com/syncedvideo/syncedvideo/http/middleware"
	"github.com/syncedvideo/syncedvideo/http/response"
)

type roomHandler struct {
	store syncedvideo.Store
	redis *redis.Client
}

func RegisterRoomHandler(router chi.Router, store syncedvideo.Store, redis *redis.Client) {
	roomHandler := newRoomHandler(store, redis)
	router.Route("/room", func(router2 chi.Router) {
		router2.Use(func(next http.Handler) http.Handler {
			return middleware.UserMiddleware(next, store.User())
		})
		router2.Post("/", roomHandler.Create)
		router2.Get("/{roomID}", roomHandler.Get)
		router2.Put("/{roomID}", roomHandler.Update)
		router2.HandleFunc("/{roomID}/connect", roomHandler.Connect)
	})
}

func newRoomHandler(s syncedvideo.Store, r *redis.Client) *roomHandler {
	return &roomHandler{
		store: s,
		redis: r,
	}
}

func (h *roomHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserCtx(r)
	room := syncedvideo.Room{OwnerUserID: user.ID}
	if err := h.store.Room().Create(&room); err != nil {
		log.Printf("error creating room: %s\n", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	log.Printf("created room id: %v\n", room.ID)
	response.WithJSON(w, room, http.StatusCreated)
}

func (h *roomHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "roomID"))
	if err != nil {
		log.Printf("error parsing uuid: %v", err)
		response.WithError(w, "room not found", http.StatusNotFound)
		return
	}
	room, err := h.store.Room().Get(id)
	if err == sql.ErrNoRows {
		response.WithError(w, "room not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("error getting room: %v", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	response.WithJSON(w, room, http.StatusOK)
}

func (h *roomHandler) Update(w http.ResponseWriter, r *http.Request) {

	panic("not implemented") // TODO: Implement
}

func (h *roomHandler) Vote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *roomHandler) Connect(w http.ResponseWriter, r *http.Request) {
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

	user := middleware.GetUserCtx(r)
	user.SetConnection(conn)
	room.Run(&user, h.store, h.redis)
	// for {
	// 	_, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		log.Printf("error reading message: %v\n", err)
	// 		break
	// 	}
	// 	log.Printf("recieved message: %v\n", msg)
	// }
}
