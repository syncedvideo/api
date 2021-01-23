package syncedvideo

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	ChatColor string    `json:"chatColor"`
	Buffering bool      `json:"buffering"`
	Time      int64     `json:"time"`
}

type UserStore interface {
	GetUser(id uuid.UUID) (User, error)
	CreateUser(u *User) (*User, error)
	UpdateUser(u *User) (*User, error)
	DeleteUser(id uuid.UUID) error
}

// NewUser returns a new user
func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		Username:  "",
		ChatColor: "",
		Buffering: false,
		Time:      0,
	}
}

// SetBuffering user prop
func (u *User) SetBuffering(buffering bool) *User {
	u.Buffering = buffering
	return u
}

// SetUsername of user
func (u *User) SetUsername(name string) *User {
	u.Username = name
	return u
}

// SetChatColor of chat
func (u *User) SetChatColor(color string) *User {
	u.ChatColor = color
	return u
}
