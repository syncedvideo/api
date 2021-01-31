package syncedvideo

import (
	"github.com/google/uuid"
)

type Store interface {
	User() UserStore
	Room() RoomStore
	Playlist() PlaylistStore
}

type UserStore interface {
	Get(id uuid.UUID) (User, error)
	Create(u *User) error
	Update(u *User) error
	Delete(id uuid.UUID) error
	// Connect(roomID uuid.UUID) error
	// Disconnect(roomID uuid.UUID) error
}

type RoomStore interface {
	Get(id uuid.UUID) (Room, error)
	Create(r *Room) error
	Update(r *Room) error
	Delete(id uuid.UUID) error
	// Connect(userID uuid.UUID) error
	// Disconnect(userID uuid.UUID) error
}

type PlaylistStore interface {
	Get(itemID uuid.UUID) (PlaylistItem, error)
	All(roomID uuid.UUID) (map[uuid.UUID]PlaylistItem, error)
	Create(p *PlaylistItem) error
	Update(p *PlaylistItem) error
	Delete(itemID uuid.UUID) error

	GetVote(itemID uuid.UUID, userID uuid.UUID) (PlaylistItemVote, error)
	CreateVote(itemID uuid.UUID, userID uuid.UUID) error
	DeleteVote(itemID uuid.UUID, userID uuid.UUID) error
}
