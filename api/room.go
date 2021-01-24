package syncedvideo

import (
	"github.com/google/uuid"
)

// Room handles state and connected users
type Room struct {
	ID            uuid.UUID           `db:"id" json:"id"`
	Name          string              `db:"name" json:"name"`
	OwnerUserID   uuid.UUID           `db:"owner_user_id"`
	Users         map[uuid.UUID]*User `json:"users"`
	Player        *Player             `json:"player"`
	Chat          *Chat               `json:"chat"`
	ConnectionHub *ConnectionHub      `json:"connectionHub"`
}

type RoomStore interface {
	GetRoom(id uuid.UUID) (Room, error)
	CreateRoom(r *Room) error
	UpdateRoom(r *Room) error
	DeleteRoom(id uuid.UUID) error
}

func NewRoom(connectionCap int) *Room {
	return &Room{
		ID:            uuid.New(),
		Player:        NewVideoPlayer(),
		Chat:          NewChat(),
		ConnectionHub: NewConnectionHub(connectionCap),
	}
}

// BroadcastSync room state with all connected users
func (r *Room) BroadcastSync() {
	r.Player.Queue.Sort()
	r.ConnectionHub.BroadcastEvent(WsEvent{
		Name: WsEventSync,
		Data: r,
	})
}

func (r *Room) BroadcastRoomSeeked(t int64) {
	r.ConnectionHub.BroadcastEvent(WsEvent{
		Name: WsEventPlayerSeeked,
		Data: t,
	})
}

func (r *Room) FindUser(id uuid.UUID) *User {
	connection, exists := r.ConnectionHub.Connections[id]
	if !exists {
		return nil
	}
	return connection.User
}
