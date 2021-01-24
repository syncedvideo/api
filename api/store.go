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
}
