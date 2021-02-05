package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/middleware"
	"github.com/syncedvideo/syncedvideo/http/request"
	"github.com/syncedvideo/syncedvideo/http/response"
)

type roomHandler struct{}

func RegisterRoomHandler(r chi.Router) {
	roomHandler := newRoomHandler()
	r.Route("/room", func(r chi.Router) {
		r.Use(middleware.UserMiddleware)
		r.Post("/", roomHandler.Create)
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(middleware.RoomMiddleware)
			r.Get("/", roomHandler.Get)
			r.Put("/", roomHandler.Update)
			r.HandleFunc("/websocket", roomHandler.WebSocket)
			r.Post("/chat", roomHandler.Chat)
		})
	})
}

func newRoomHandler() *roomHandler {
	return &roomHandler{}
}

func (h *roomHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := request.GetUserCtx(r)
	room := syncedvideo.Room{OwnerUserID: user.ID}
	if err := syncedvideo.Config.Store.Room().Create(&room); err != nil {
		log.Printf("error creating room: %s\n", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	log.Printf("created room id: %v\n", room.ID)
	response.WithJSON(w, room, http.StatusCreated)
}

func (h *roomHandler) Get(w http.ResponseWriter, r *http.Request) {
	room := request.GetRoomCtx(r)
	err := syncedvideo.Config.Store.Room().WithUsers(&room)
	if err != nil {
		fmt.Println(err)
	}
	response.WithJSON(w, room, http.StatusOK)
}

func (h *roomHandler) Update(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *roomHandler) Vote(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h *roomHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
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
	user.ConnectionID = uuid.New()
	user.Connection = conn

	defer func() {
		user.Connection.Close()
		syncedvideo.Config.Store.Room().Leave(&room, &user)
		room.Publish(syncedvideo.WebSocketMessageLeave, user)
	}()
	syncedvideo.Config.Store.Room().Join(&room, &user)
	room.Publish(syncedvideo.WebSocketMessageJoin, user)
	room.Run(&user)
}

type ChatData struct {
	Message string `json:"message"`
}

func (h *roomHandler) Chat(w http.ResponseWriter, r *http.Request) {
	data := ChatData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("error decoding data: %v\n", err)
		response.WithError(w, "something went wrong", http.StatusBadRequest)
		return
	}
	if data.Message == "" {
		response.WithError(w, "message is required", http.StatusBadRequest)
		return
	}
	chatMessage := syncedvideo.NewChatMessage(request.GetUserCtx(r), data.Message)
	room := request.GetRoomCtx(r)
	room.Publish(syncedvideo.WebSocketMessageChat, chatMessage)
}
