package syncedvideo

import (
	"github.com/google/uuid"
)

// Room handles state and connected users
type Room struct {
	ID            uuid.UUID                  `db:"id" json:"id"`
	Name          string                     `db:"name" json:"name"`
	Description   string                     `db:"description" json:"description"`
	OwnerUserID   uuid.UUID                  `json:"ownerUserId" db:"owner_user_id"`
	Users         map[uuid.UUID]*User        `json:"users"`
	Player        *Player                    `json:"player"`
	Chat          *Chat                      `json:"chat"`
	ConnectionHub *ConnectionHub             `json:"connectionHub"`
	PlaylistItems map[uuid.UUID]PlaylistItem `json:"playlistItems"`
}

type PlaylistItem struct {
	ID     uuid.UUID `db:"id"`
	RoomID uuid.UUID `db:"room_id"`
	UserID uuid.UUID `db:"user_id"`
	Votes  []PlaylistItemVote
}

type PlaylistItemVote struct {
	ID     uuid.UUID `db:"id"`
	ItemID uuid.UUID `db:"item_id"`
	UserID uuid.UUID `db:"user_id"`
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
