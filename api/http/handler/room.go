package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/middleware"
	"github.com/syncedvideo/syncedvideo/http/request"
	"github.com/syncedvideo/syncedvideo/http/response"
	"github.com/syncedvideo/syncedvideo/youtube"
)

type roomHandler struct {
	YouTubeAPIKey string
}

func RegisterRoomHandler(r chi.Router, ytAPIKey string) {
	roomHandler := newRoomHandler(ytAPIKey)
	r.Route("/room", func(r chi.Router) {
		r.Use(middleware.UserMiddleware)
		r.Post("/", roomHandler.Create)
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(middleware.RoomMiddleware)
			r.Get("/", roomHandler.Get)
			r.Put("/", roomHandler.Update)
			r.HandleFunc("/websocket", roomHandler.WebSocket)
			r.Post("/chat", roomHandler.Chat)
			r.Get("/video-info", roomHandler.VideoInfo)
		})
	})
}

func newRoomHandler(ytAPIKey string) *roomHandler {
	return &roomHandler{YouTubeAPIKey: ytAPIKey}
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
		response.WithError(w, "missing message", http.StatusBadRequest)
		return
	}
	chatMessage := syncedvideo.NewChatMessage(request.GetUserCtx(r), data.Message)
	room := request.GetRoomCtx(r)
	room.Publish(syncedvideo.WebSocketMessageChat, chatMessage)
}

func (h *roomHandler) VideoInfo(w http.ResponseWriter, r *http.Request) {
	videoURL := r.URL.Query().Get("url")
	if videoURL == "" {
		response.WithError(w, "missing url query param", http.StatusBadRequest)
		return
	}

	videoID := youtube.ExtractVideoID(videoURL)
	if videoID == "" {
		response.WithError(w, "missing video ID", http.StatusBadRequest)
		return
	}

	// try to get video from cache
	video, err := youtube.GetVideoFromCache(syncedvideo.Config.Redis, videoID)
	if err != nil && err != redis.Nil {
		log.Printf("GetVideoFromCache failed: %v\n", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	if video.ID != "" {
		response.WithJSON(w, syncedvideo.NewVideo(video), http.StatusOK)
		return
	}

	// try to get video from YouTube API
	client := youtube.New(h.YouTubeAPIKey)
	video, err = client.GetVideo(videoID)
	if err != nil && err != youtube.ErrNoResults {
		log.Printf("GetVideo failed: %v\n", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	if err == youtube.ErrNoResults {
		log.Printf("no cache key for %v\n", videoID)
	} else {
		log.Printf("got video from cache: %v\n", video)
	}

	// add video to cache
	err = youtube.CacheVideo(syncedvideo.Config.Redis, video)
	if err != nil {
		log.Printf("CacheVideo failed: %v\n", err)
		response.WithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	log.Printf("added cache key: %v\n", videoID)

	if video.ID == "" {
		response.WithJSON(w, nil, http.StatusOK)
		return
	}
	response.WithJSON(w, syncedvideo.NewVideo(video), http.StatusOK)
}
