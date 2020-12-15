package room

import (
	"github.com/google/uuid"
)

// Room handles state and connected users
type Room struct {
	ID            uuid.UUID           `json:"id"`
	Users         map[uuid.UUID]*User `json:"users"`
	Player        *Player             `json:"player"`
	Chat          *Chat               `json:"chat"`
	ConnectionHub *ConnectionHub      `json:"connectionHub"`
}

// NewRoom returns a new room
func NewRoom(connectionCap int) *Room {
	return &Room{
		ID:            uuid.New(),
		Player:        NewVideoPlayer(),
		Chat:          NewChat(),
		ConnectionHub: NewConnectionHub(connectionCap),
	}
}

// Sync room state with all connected users
func (r *Room) Sync() {
	r.ConnectionHub.BroadcastEvent(WsEvent{
		Name: WsEventSync,
		Data: r,
	})
}
func (r *Room) FindUser(id uuid.UUID) *User {
	connection, exists := r.ConnectionHub.Connections[id]
	if !exists {
		return nil
	}
	return connection.User
}
