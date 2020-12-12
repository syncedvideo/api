package room

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	ChatColor string    `json:"chatColor"`
}

// NewUser returns a new user
func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		Username:  "",
		ChatColor: "",
	}
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
