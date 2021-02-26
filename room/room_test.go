package room

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRun(t *testing.T) {
	room := newRoom()
	go room.Run()
	defer room.Close()

	t.Run("join and leave", func(t *testing.T) {
		users := []userMock{newUser(), newUser(), newUser()}
		for _, user := range users {
			room.userJoin <- user
		}
		if len(room.users) != len(users) {
			t.Errorf("join failed: got %v want %v", len(room.users), len(users))
		}

		for _, user := range users {
			room.userLeave <- user
		}
		if len(room.users) != 0 {
			t.Errorf("leave failed: got %v want 0", len(room.users))
		}
	})

	t.Run("add and remove video to playlist", func(t *testing.T) {
		video := newVideo()
		room.playlistAdd <- video
		_, found := room.playlist[video]
		if !found {
			t.Error("video was not added to playlist")
		}
		room.playlistRemove <- video
		_, found = room.playlist[video]
		if found {
			t.Error("video was not removed from playlist")
		}
	})
}

func newRoom() Room {
	return New()
}

type userMock struct {
	id uuid.UUID
}

func (user userMock) ID() fmt.Stringer {
	return user.id
}

func newUser() userMock {
	return userMock{id: uuid.New()}
}

type videoMock struct {
	id uuid.UUID
}

func (v videoMock) ID() fmt.Stringer {
	return v.id
}

func (v videoMock) Provider() string {
	return "YouTube"
}

func (v videoMock) ProviderID() string {
	return uuid.NewString()
}

func (v videoMock) Title() string {
	return "Test Title"
}

func (v videoMock) Author() string {
	return "Test Author"
}

func (v videoMock) Duration() time.Duration {
	return time.Second * time.Duration(rand.Int())
}

func (v videoMock) Views() int {
	return rand.Int()
}

func (v videoMock) Likes() int {
	return rand.Int()
}

func (v videoMock) Dislikes() int {
	return rand.Int()
}

func newVideo() videoMock {
	return videoMock{id: uuid.New()}
}
