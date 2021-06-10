package syncedvideo

import (
	"time"

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
			Current: nil,
		},
		eventManager: eventManager,
	}
}

func (r *Room) Play(video *Video) {
	r.VideoPlayer.Current = video
	r.VideoPlayer.StartedAt = time.Now()
	r.VideoPlayer.PausedAt = time.Time{}
	event := NewEvent(EventPlay, video)
	r.eventManager.Publish(r.ID, event)
}

func (r *Room) SendChatMessage(chatMessage ChatMessage) {
	event := NewEvent(EventChat, chatMessage)
	r.eventManager.Publish(r.ID, event)
}

type VideoPlayer struct {
	Current   *Video
	StartedAt time.Time
	PausedAt  time.Time
}

type Video struct {
	ID         string `json:"id"`
	Provider   string `json:"provider"`
	ProviderID string `json:"providerId"`
	Title      string `json:"title"`
}

func (p *VideoPlayer) Playing() bool {
	return p.StartedAt.After(p.PausedAt)
}
