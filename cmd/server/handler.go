package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/syncedvideo/backend/room"
)

// WsActionHandler manages WebSocket action handler methods
type WsActionHandler struct {
	WsAction *room.WsAction
	Room     *room.Room
	User     *room.User
}

// NewWsActionHandler returns a new WsActionHandler
func NewWsActionHandler(a *room.WsAction, r *room.Room, u *room.User) *WsActionHandler {
	return &WsActionHandler{
		WsAction: a,
		Room:     r,
		User:     u,
	}
}

// Handle incoming WebSocket actions
func (handler *WsActionHandler) Handle() {
	switch handler.WsAction.Name {

	// User actions
	case room.WsActionUserSetUsername:
		handler.handleUserSetUsername()
	case room.WsActionUserSetColor:
		handler.handleUserSetColor()

	// Player actions
	case room.WsActionPlayerInit:
		handler.handlePlayerInit()
	case room.WsActionPlayerSkip:
		handler.handlePlayerSkip()
	case room.WsActionPlayerTogglePlaying:
		handler.handlePlayerTogglePlaying()

	// Queue actions
	case room.WsActionQueueAdd:
		handler.handleQueueAdd()
	case room.WsActionQueueRemove:
		handler.handleQueueRemove()
	case room.WsActionQueueVote:
		handler.handleQueueVote()

	// Chat actions
	case room.WsActionChatMessage:
		handler.handleChatMessage()
	}

	// Sync room state after handling the action
	handler.Room.Sync()
}

func (handler *WsActionHandler) handleUserSetUsername() {
	username := ""
	err := json.Unmarshal(handler.WsAction.Data, &username)
	if err != nil {
		log.Println("ClientSetUsername error:", err)
		return
	}
	handler.User.SetUsername(username)
}

func (handler *WsActionHandler) handleUserSetColor() {
	color := ""
	err := json.Unmarshal(handler.WsAction.Data, &color)
	if err != nil {
		log.Println("ClientSetUsername error:", err)
		return
	}
	handler.User.SetChatColor(color)
}

func (handler *WsActionHandler) handlePlayerInit() {
	log.Println("TODO handlePlayerInit")
}

func (handler *WsActionHandler) handlePlayerSkip() {
	if len(handler.Room.VideoPlayer.Queue.Videos) >= 1 {
		// Set current video
		handler.Room.VideoPlayer.CurrentVideo = handler.Room.VideoPlayer.Queue.Videos[0]
		// Remove current video from queue
		handler.Room.VideoPlayer.Queue.Remove(handler.Room.VideoPlayer.CurrentVideo.ID)

		log.Panicln("handlePlayerSkip: Video skipped by user:", handler.User)
		return
	}
	log.Println("handlePlayerSkip: Queue is empty")
}

func (handler *WsActionHandler) handlePlayerTogglePlaying() {
	if handler.Room.VideoPlayer.CurrentVideo == nil {
		log.Println("handlePlayerTogglePlaying: CurrentVideo is nil")
		return
	}
	handler.Room.VideoPlayer.Playing = !handler.Room.VideoPlayer.Playing
}

func (handler *WsActionHandler) handleQueueAdd() {
	video := &room.Video{}
	err := json.Unmarshal(handler.WsAction.Data, video)
	if err != nil {
		log.Println("handleQueueAdd error:", err)
		return
	}
	handler.Room.VideoPlayer.Play(video)
}

func (handler *WsActionHandler) handleQueueRemove() {
	videoUUIDString := ""
	err := json.Unmarshal(handler.WsAction.Data, &videoUUIDString)
	if err != nil {
		log.Println("handleQueueRemove error:", err)
		return
	}
	videoUUID, _ := uuid.Parse(videoUUIDString)
	if videoUUID.String() != "" {
		handler.Room.VideoPlayer.Queue.Remove(videoUUID)
	}
}

func (handler *WsActionHandler) handleQueueVote() {
	videoUUIDString := ""
	err := json.Unmarshal(handler.WsAction.Data, &videoUUIDString)
	if err != nil {
		log.Println("handleQueueVote error:", err)
		return
	}

	videoUUID, err := uuid.Parse(videoUUIDString)
	if err != nil {
		log.Println("handleQueueVote error:", err)
		return
	}

	video := handler.Room.VideoPlayer.Queue.Find(videoUUID)
	if video == nil {
		log.Println("handleQueueVote: Video %w not found", videoUUID)
		return
	}

	handler.Room.VideoPlayer.Queue.ToggleVote(handler.User, video)
}

func (handler *WsActionHandler) handleChatMessage() {
	text := ""
	err := json.Unmarshal(handler.WsAction.Data, &text)
	if err != nil {
		log.Println("ChatMessage error:", err)
		return
	}
	handler.Room.Chat.NewMessage(handler.User, text)
}
