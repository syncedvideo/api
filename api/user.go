package syncedvideo

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"username"`
	Color     string    `db:"color" json:"chatColor"`
	IsAdmin   bool      `db:"is_admin" json:"isAdmin"`
	Buffering bool      `json:"-"`
	Time      int64     `json:"-"`
}

func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		Name:      "",
		Color:     "",
		Buffering: false,
		Time:      0,
	}
}

func (u *User) SetBuffering(buffering bool) *User {
	u.Buffering = buffering
	return u
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
