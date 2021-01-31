package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/middleware"
	"github.com/syncedvideo/syncedvideo/http/request"
	"github.com/syncedvideo/syncedvideo/http/response"
)

type roomHandler struct {
	store syncedvideo.Store
	redis *redis.Client
}

func RegisterRoomHandler(r chi.Router, s syncedvideo.Store, rc *redis.Client) {
	roomHandler := newRoomHandler(s, rc)
	r.Route("/room", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middleware.UserMiddleware(next, s.User())
		})
		r.Post("/", roomHandler.Create)
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return middleware.RoomMiddleware(next, s.Room())
			})
			r.Get("/", roomHandler.Get)
			r.Put("/", roomHandler.Update)
			r.HandleFunc("/connect", roomHandler.Connect)
			r.Post("/chat", roomHandler.Chat)
		})
	})
}

func newRoomHandler(store syncedvideo.Store, redis *redis.Client) *roomHandler {
	return &roomHandler{
		store,
		redis,
	}
}

func (h *roomHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserCtx(r)
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
	response.WithJSON(w, request.GetRoomCtx(r), http.StatusOK)
}

func (h *roomHandler) Update(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *roomHandler) Vote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *roomHandler) Connect(w http.ResponseWriter, r *http.Request) {
	room := request.GetRoomCtx(r)
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

	user := request.GetUserCtx(r)
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

func (h *roomHandler) Chat(w http.ResponseWriter, r *http.Request) {
	log.Println("hi from chat")
}
