package syncedvideo

import "github.com/google/uuid"

type Store interface {
	User() UserStore
	Room() RoomStore
}

type UserStore interface {
	Get(id uuid.UUID) (User, error)
	Create(u *User) error
	Update(u *User) error
	Delete(id uuid.UUID) error
}

type RoomStore interface {
	Get(id uuid.UUID) (Room, error)
	Create(r *Room) error
	Update(r *Room) error
	Delete(id uuid.UUID) error
	GetPlaylistItem(r *Room, id uuid.UUID) (PlaylistItem, error)
	GetAllPlaylistItems(roomID uuid.UUID) (map[uuid.UUID]PlaylistItem, error)
	CreatePlaylistItem(r *Room, p *PlaylistItem) error
	UpdatePlaylistItem(r *Room, p *PlaylistItem) error
	DeletePlaylistItem(r *Room, id uuid.UUID) error
}
