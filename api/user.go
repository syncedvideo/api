package syncedvideo

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"username"`
	Color     string    `db:"color" json:"chatColor"`
	IsAdmin   bool      `db:"is_admin" json:"isAdmin"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"UpdatedAt"`

	Buffering bool  `json:"buffering"`
	Time      int64 `json:"time"`
}

type UserStore interface {
	GetUser(id uuid.UUID) (User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id uuid.UUID) error
}

// NewUser returns a new user
func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		Name:      "",
		Color:     "",
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
	u.Name = name
	return u
}

// SetChatColor of chat
func (u *User) SetChatColor(color string) *User {
	u.Color = color
	return u
}
