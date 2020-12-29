package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	roomPackage "github.com/syncedvideo/backend/room"
)

var addr = flag.String("addr", ":3000", "http service address")
var frontendURL = flag.String("frontendURL", "localhost:8080", "url of frontend")

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
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

const connectionCap = 10

func postRoomHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	newRoom := roomPackage.NewRoom(connectionCap)
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

	// Check userID cookie
	userIDCookie, err := r.Cookie("userID")
	var userIDCookieUUID uuid.UUID
	if err == nil {
		log.Println("cookie:", userIDCookie.Value)
		userIDCookieUUID, _ = uuid.Parse(userIDCookie.Value)
	}

	user := room.FindUser(userIDCookieUUID)
	if user == nil {
		// create user
		user = roomPackage.NewUser()
	}

	// Create userID cookie header
	userIDCookie = &http.Cookie{
		Name:     "userID",
		Value:    user.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour * 360),
	}
	header := make(http.Header)
	header["Set-Cookie"] = []string{userIDCookie.String()}

	// upgrade http connection to websocket connection with userID cookie
	wsConn, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		log.Print("upgrader.Upgrade error:", err)
		return
	}

	// connect user to room
	_, err = room.ConnectionHub.Connect(user, wsConn)
	if err != nil {
		log.Println(err)
		return
	}
	room.BroadcastSync()

	// handle disconnect
	defer func() {
		room.ConnectionHub.Disconnect(user, wsConn)
		wsConn.Close()
		room.BroadcastSync()
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

		NewWsActionHandler(wsAction, room, user).Handle()
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

	videoSearch, err := roomPackage.NewVideoSearch(os.Getenv("YOUTUBE_API_KEY")).Do(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("youTubeSearchHandler error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&videoSearch)
}
