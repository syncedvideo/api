package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	roomPackage "github.com/syncedvideo/backend/room"
	"github.com/syncedvideo/backend/youtube"
)

var addr = flag.String("addr", "localhost:3000", "http service address")
var frontendURL = flag.String("frontendURL", "localhost:8080", "url of frontend")

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	flag.Parse()
	router := chi.NewRouter()
	router.Post("/room", postRoomHandler)
	router.HandleFunc("/room/{roomID}", roomWebSocketHandler)
	router.Get("/search/youtube", searchYouTubeHandler)
	http.ListenAndServe(*addr, router)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

var rooms = make(map[uuid.UUID]*roomPackage.Room)

func postRoomHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	newRoom := roomPackage.NewRoom()
	rooms[newRoom.ID] = newRoom
	log.Println("Created new room", newRoom.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRoom)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func roomWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// find room
	roomIDParam := chi.URLParam(r, "roomID")
	if roomIDParam == "" {
		log.Println("roomIDParam is empty")
		return
	}
	roomID, err := uuid.Parse(roomIDParam)
	if err != nil {
		log.Println("uuid.Parse error:", err)
		return
	}
	room, found := rooms[roomID]
	if !found {
		log.Println("room not found:", roomID.String())
		return
	}

	// create user
	user := roomPackage.NewUser()

	// upgrade http connection to websocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrader.Upgrade error:", err)
		return
	}

	// connect user to room
	room.ConnectionHub.Connect(user, wsConn)
	room.Sync()

	// handle disconnect
	defer func() {
		room.ConnectionHub.Disconnect(user)
		wsConn.Close()
	}()

	// handle incoming messages
	for {
		_, msg, err := wsConn.ReadMessage()

		if err != nil {
			log.Println("ReadMessage:", err)
			break
		}

		wsAction := &roomPackage.WsAction{}
		err = json.Unmarshal(msg, wsAction)
		if err != nil {
			log.Println(err)
			continue
		}

		wsActionHandler := NewWsActionHandler(wsAction, room, user)
		wsActionHandler.Handle()
	}
}

func searchYouTubeHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	query := r.URL.Query().Get("query")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("youTubeSearchHandler: query is empty")
		return
	}

	yt := youtube.New(os.Getenv("YOUTUBE_API_KEY"))
	result, err := yt.VideoSearch(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("youTubeSearchHandler error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
