package main

import (
	"encoding/json"
	"log"
	"time"

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
	case room.WsActionUserSetBuffering:
		handler.handleUserSetBuffering()
	case room.WsActionUserSetUsername:
		handler.handleUserSetUsername()
	case room.WsActionUserSetColor:
		handler.handleUserSetColor()

	// Player actions
	case room.WsActionPlayerPlay:
		handler.handlePlayerPlay()
	case room.WsActionPlayerPause:
		handler.handlePlayerPause()
	case room.WsActionPlayerSkip:
		handler.handlePlayerSkip()
	case room.WsActionPlayerSeek:
		handler.handlePlayerSeek()

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
	handler.Room.BroadcastSync()
}

func (handler *WsActionHandler) handleUserSetBuffering() {
	var buffering bool
	err := json.Unmarshal(handler.WsAction.Data, &buffering)
	if err != nil {
		log.Println("handleUserSetBuffering error:", err)
		return
	}
	handler.User.SetBuffering(buffering)
}

func (handler *WsActionHandler) handleUserSetUsername() {
	username := ""
	err := json.Unmarshal(handler.WsAction.Data, &username)
	if err != nil {
		log.Println("handleUserSetUsername error:", err)
		return
	}
	handler.User.SetUsername(username)
}

func (handler *WsActionHandler) handleUserSetColor() {
	color := ""
	err := json.Unmarshal(handler.WsAction.Data, &color)
	if err != nil {
		log.Println("handleUserSetColor error:", err)
		return
	}
	handler.User.SetChatColor(color)
}

func (handler *WsActionHandler) handlePlayerPlay() {
	if handler.Room.Player.Video == nil {
		log.Println("handlePlayerTogglePlaying: Video is nil")
		return
	}
	handler.Room.Player.Playing = true
}

func (handler *WsActionHandler) handlePlayerPause() {
	if handler.Room.Player.Video == nil {
		log.Println("handlePlayerTogglePlaying: Video is nil")
		return
	}
	handler.Room.Player.Playing = false
}

func (handler *WsActionHandler) handlePlayerSkip() {
	if len(handler.Room.Player.Queue.Videos) >= 1 {
		// Set current video
		handler.Room.Player.Video = handler.Room.Player.Queue.Videos[0]
		// Remove current video from queue
		handler.Room.Player.Queue.Remove(handler.Room.Player.Video.ID)

		log.Println("handlePlayerSkip: Video skipped by user:", handler.User)
		return
	}
	log.Println("handlePlayerSkip: Queue is empty")
}

func (handler *WsActionHandler) handlePlayerSeek() {
	var t time.Duration
	err := json.Unmarshal(handler.WsAction.Data, &t)
	if err != nil {
		log.Println("e error:", err)
		return
	}
	handler.Room.BroadcastRoomSeeked(t)
}

func (handler *WsActionHandler) handleQueueAdd() {
	video := &room.Video{}
	err := json.Unmarshal(handler.WsAction.Data, video)
	if err != nil {
		log.Println("handleQueueAdd error:", err)
		return
	}
	log.Println("handleQueueAdd:", video)
	if handler.Room.Player.Video == nil {
		handler.Room.Player.Play(video)
		return
	}
	video.AddVote(handler.User)
	handler.Room.Player.Queue.Add(handler.User, video)
}

func (handler *WsActionHandler) handleQueueRemove() {
	idString := ""
	err := json.Unmarshal(handler.WsAction.Data, &idString)
	if err != nil {
		log.Println("handleQueueRemove error:", err)
		return
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println("handleQueueRemove error:", err)
		return
	}
	handler.Room.Player.Queue.Remove(id)
}

func (handler *WsActionHandler) handleQueueVote() {
	idString := ""
	err := json.Unmarshal(handler.WsAction.Data, &idString)
	if err != nil {
		log.Println("handleQueueVote error:", err)
		return
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		log.Println("handleQueueVote error:", err)
		return
	}

	video := handler.Room.Player.Queue.Find(id)
	if video == nil {
		log.Println("handleQueueVote: Video %w not found", id)
		return
	}

	handler.Room.Player.Queue.ToggleVote(handler.User, video)
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
