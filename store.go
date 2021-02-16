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
	Get(userID uuid.UUID) (User, error)
	Create(u *User) error
	Update(u *User) error
	Delete(userID uuid.UUID) error
	GetCurrentRooms(userID uuid.UUID) ([]Room, error)
}

type RoomStore interface {
	Get(roomID uuid.UUID) (Room, error)
	Create(r *Room) error
	Update(r *Room) error
	Delete(roomID uuid.UUID) error
	Join(r *Room, u *User) error
	Leave(r *Room, u *User) error
	GetUsers(r *Room) ([]User, error)
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
