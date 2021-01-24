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
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/postgres"
)

var (
	apiHTTPPort         = os.Getenv("API_HTTP_PORT")
	apiPostgresHost     = os.Getenv("API_POSTGRES_HOST")
	apiPostgresPort     = os.Getenv("API_POSTGRES_PORT")
	apiPostgresDB       = os.Getenv("API_POSTGRES_DB")
	apiPostgresUser     = os.Getenv("API_POSTGRES_USER")
	apiPostgresPassword = os.Getenv("API_POSTGRES_PASSWORD")
	apiRedisHost        = os.Getenv("API_REDIS_HOST")
	apiRedisPort        = os.Getenv("API_REDIS_PORT")
)

func main() {
	flag.Parse()

	// init store
	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", apiPostgresHost, apiPostgresUser, apiPostgresPassword, apiPostgresDB)
	store, err := postgres.NewStore(postgresDsn)
	if err != nil {
		panic(err)
	}

	r := syncedvideo.Room{ID: uuid.New(), Name: "test"}
	err = store.CreateRoom(&r)
	if err != nil {
		panic(err)
	}

	// init redis client
	redisOpts, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s", apiRedisHost, apiRedisPort))
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisOpts)
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	// init http server
	h := RegisterHandlers(store, redisClient)
	log.Printf("http server listening on port %s\n", apiHTTPPort)
	go http.ListenAndServe(fmt.Sprintf(":%s", apiHTTPPort), h)

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

var rooms = make(map[uuid.UUID]*syncedvideo.Room)

const connectionCap = 10

func postRoomHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	newRoom := syncedvideo.NewRoom(connectionCap)
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
		user = syncedvideo.NewUser()
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

		wsAction := &syncedvideo.WsAction{}
		err = json.Unmarshal(msg, wsAction)
		if err != nil {
			log.Println(err)
			continue
		}

		NewWsActionHandler(wsAction, room, user).Handle()
	}
}
