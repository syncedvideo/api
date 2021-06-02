package syncedvideo

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomStore interface {
	CreateRoom(room *Room)
	GetRoom(id string) Room
}

type Server struct {
	store RoomStore
	http.Handler
	pubSub PubSub
}

const (
	jsonContentType = "application/json"
	userCookieName  = "user"
)

func NewServer(store RoomStore, pubSub PubSub) *Server {
	server := new(Server)
	server.store = store
	server.pubSub = pubSub

	router := chi.NewMux()
	router.Post("/rooms", server.postRoomHandler)
	router.Get("/rooms/{id}", server.getRoomHandler)
	router.Get("/rooms/{id}/ws", server.webSocket)
	router.Post("/rooms/{id}/chat", server.postChatHandler)
	server.Handler = router

	return server
}

func (s *Server) postRoomHandler(w http.ResponseWriter, r *http.Request) {
	room, _ := NewRoom(r.Body)
	s.store.CreateRoom(&room)

	http.SetCookie(w, NewUserCookie())
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(room)
}

func (s *Server) getRoomHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	room := s.store.GetRoom(id)
	if room.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.SetCookie(w, NewUserCookie())
	w.Header().Set("Content-Type", jsonContentType)

	json.NewEncoder(w).Encode(room)
}

func (s *Server) postChatHandler(w http.ResponseWriter, r *http.Request) {

	_, err := r.Cookie(userCookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	roomID := chi.URLParam(r, "id")
	room := s.store.GetRoom(roomID)

	if room.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusCreated)

	bodyData := ChatMessage{}
	json.NewDecoder(r.Body).Decode(&bodyData)
	chatMsg := NewChatMessage(bodyData.Author, bodyData.Message)

	chatMsgBytes, _ := json.Marshal(chatMsg)

	event := NewEvent(EventChat, chatMsgBytes)
	go s.pubSub.Publish(room.ID, event)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) webSocket(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	room := s.store.GetRoom(roomID)
	if room.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch := s.pubSub.Subscribe(roomID)

	for event := range ch {
		log.Printf("received event: %v\n", event)

		err := conn.WriteJSON(event)
		if err != nil {
			log.Printf("error writing json: %v\n", err)
			break
		}
	}
}

func NewUserCookie() *http.Cookie {
	return &http.Cookie{
		Name:     userCookieName,
		HttpOnly: false,
		Value:    uuid.NewString(),
		Expires:  time.Now().Add(time.Hour * 24 * 360),
	}
}
