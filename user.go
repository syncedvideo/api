package syncedvideo

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/syncedvideo/syncedvideo/sillyname"
)

type User struct {
	ID           uuid.UUID       `db:"id" json:"id"`
	Name         string          `db:"name" json:"name"`
	Color        string          `db:"color" json:"color"`
	IsAdmin      bool            `db:"is_admin" json:"isAdmin"`
	Connection   *websocket.Conn `json:"-"`
	ConnectionID uuid.UUID       `json:"-"`
}

func NewUser() User {
	return User{
		ID:      uuid.New(),
		Name:    sillyname.New(),
		IsAdmin: false,
	}
}

func (u *User) SetUsername(name string) *User {
	u.Name = name
	return u
}

func (u *User) SetChatColor(color string) *User {
	u.Color = color
	return u
}

func (u *User) CanUpdateRoom() bool {
	return true
}
