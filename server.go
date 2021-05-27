package syncedvideo

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type RoomStore interface {
	GetRoom(id string) Room
}

type RoomPubSub interface {
	Publish(roomID string, event RoomEvent)
	Subscribe(roomID string) <-chan RoomEvent
}

type Server struct {
	store RoomStore
	http.Handler
	pubSub RoomPubSub
}

const jsonContentType = "application/json"

func NewServer(store RoomStore, pubSub RoomPubSub) *Server {
	server := new(Server)
	server.store = store
	server.pubSub = pubSub

	router := chi.NewMux()
	router.Get("/rooms/{id}", server.getRoomHandler)
	router.Get("/rooms/{id}/ws", server.webSocket)
	router.Post("/rooms/{id}/chat", server.postChatHandler)
	server.Handler = router

	return server
}

func (s *Server) getRoomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", jsonContentType)
	id := chi.URLParam(r, "id")

	room := s.store.GetRoom(id)
	if room.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(room)
}

func (s *Server) postChatHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")
	room := s.store.GetRoom(roomID)

	if room.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(http.StatusCreated)

	chatMsg1 := ChatMessage{}
	json.NewDecoder(r.Body).Decode(&chatMsg1)
	chatMsg := NewChatMessage(chatMsg1.Author, chatMsg1.Message)

	chatMsgBytes, _ := json.Marshal(chatMsg)

	event := NewRoomEvent(1, chatMsgBytes)
	s.pubSub.Publish(room.ID, event)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) webSocket(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

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
