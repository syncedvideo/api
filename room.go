package syncedvideo

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Room struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Video        *Video   `json:"video"`
	Playlist     []*Video `json:"playlist"`
	eventManager EventManager
}

type Video struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewRoom(eventManager EventManager) Room {
	return Room{
		ID:           uuid.NewString(),
		eventManager: eventManager,
	}
}

func (r *Room) PlayVideo(video *Video) {
	r.Video = video
	videoB, _ := json.Marshal(video)
	event := NewEvent(EventPlayVideo, videoB)
	r.eventManager.Publish(r.ID, event)
}

func (r *Room) SendChatMessage(chatMessage ChatMessage) {
	chatMessageB, _ := json.Marshal(&chatMessage)
	event := NewEvent(EventChat, chatMessageB)
	r.eventManager.Publish(r.ID, event)
}
