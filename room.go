package syncedvideo

import (
	"github.com/google/uuid"
)

type Room struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	VideoPlayer  *VideoPlayer `json:"videoPlayer"`
	eventManager EventManager
}

func NewRoom(eventManager EventManager) Room {
	return Room{
		ID: uuid.NewString(),
		VideoPlayer: &VideoPlayer{
			CurrentVideo: nil,
		},
		eventManager: eventManager,
	}
}

func (r *Room) SendChatMessage(chatMessage ChatMessage) {
	event := NewEvent(EventChat, chatMessage)
	r.eventManager.Publish(r.ID, event)
}
