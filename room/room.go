package room

import (
	"github.com/google/uuid"
)

// Room handles state and connected clients
type Room struct {
	ID            uuid.UUID           `json:"id"`
	Users         map[uuid.UUID]*User `json:"users"`
	VideoPlayer   *VideoPlayer        `json:"videoPlayer"`
	Chat          *Chat               `json:"chat"`
	ConnectionHub *ConnectionHub      `json:"-"`
}

// NewRoom returns a new room
func NewRoom() *Room {
	return &Room{
		ID:            uuid.New(),
		Users:         make(map[uuid.UUID]*User),
		VideoPlayer:   NewVideoPlayer(),
		Chat:          NewChat(),
		ConnectionHub: NewConnectionHub(),
	}
}

// Sync room state with all connected users
func (r *Room) Sync() {
	r.ConnectionHub.BroadcastEvent(WsEvent{
		Name: WsEventSync,
		Data: r,
	})
}
