package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/syncedvideo/syncedvideo/postgres"
	roomPackage "github.com/syncedvideo/syncedvideo/room"
	"github.com/syncedvideo/syncedvideo/youtube"
)

var (
	postgresHost     = os.Getenv("APP_POSTGRES_HOST")
	postgresPort     = os.Getenv("APP_POSTGRES_PORT")
	postgresDB       = os.Getenv("APP_POSTGRES_DB")
	postgresUser     = os.Getenv("APP_POSTGRES_USER")
	postgresPassword = os.Getenv("APP_POSTGRES_PASSWORD")
	redisHost        = os.Getenv("APP_REDIS_HOST")
	redisPort        = os.Getenv("APP_REDIS_PORT")
	youTubeAPIKey    = os.Getenv("APP_YOUTUBE_API_KEY")
)

var addr = flag.String("addr", ":3000", "http service address")

func main() {
	flag.Parse()

	video, err := youtube.GetVideoInfo("6d6L4-ADF-M")
	log.Println(video, err)

	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresUser, postgresPassword, postgresDB)
	store, err := postgres.NewStore(postgresDsn)
	if err != nil {
		panic(err)
	}

	redisAddr := fmt.Sprintf("redis://%s:%s", redisHost, redisPort)
	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	h := RegisterHandlers(store, redisClient)
	log.Printf("http server listening on port %s\n", *addr)
	go http.ListenAndServe(*addr, h)

	runtime.Goexit()
}

// pubsub := redisClient.Subscribe(context.Background(), "test")
// 	ch := pubsub.Channel()

// 	go func() {
// 		for msg := range ch {
// 			fmt.Println("received ", msg.Payload)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			redisClient.Publish(context.Background(), "test", time.Now().String())
// 			time.Sleep(time.Millisecond * 1000)
// 		}
// 	}()

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

	// ping to keep connection alive
	go func() {
		for {
			time.Sleep(time.Second * 5)
			wsConn.WriteMessage(websocket.TextMessage, []byte("ping"))
		}
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

	videoSearch, err := roomPackage.NewVideoSearch(youTubeAPIKey).Do(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("youTubeSearchHandler error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&videoSearch)
}

// newRedisClient returns a Redis client
func newRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
