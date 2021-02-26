package room

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID             fmt.Stringer `json:"id" db:"id"`
	users          map[User]bool
	playlist       map[Video]bool
	userJoin       chan User
	userLeave      chan User
	playlistAdd    chan Video
	playlistRemove chan Video
}

type User interface {
	ID() fmt.Stringer
}

type Video interface {
	ID() fmt.Stringer
	Provider() string
	ProviderID() string
	Title() string
	Author() string
	Duration() time.Duration
	Views() int
	Likes() int
	Dislikes() int
}

func New() Room {
	return Room{
		ID:             uuid.New(),
		users:          make(map[User]bool),
		playlist:       make(map[Video]bool),
		userJoin:       make(chan User),
		userLeave:      make(chan User),
		playlistAdd:    make(chan Video),
		playlistRemove: make(chan Video),
	}
}

func (r *Room) Run() {
	for {
		select {
		case user := <-r.userJoin:
			r.users[user] = true
		case user := <-r.userLeave:
			delete(r.users, user)
		case video := <-r.playlistAdd:
			r.playlist[video] = true
		case video := <-r.playlistRemove:
			delete(r.playlist, video)
		default:
		}
	}
}

func (r *Room) Close() {
	close(r.userJoin)
	close(r.userLeave)
	close(r.playlistAdd)
	close(r.playlistRemove)
}
