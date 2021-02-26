package room

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestRun(t *testing.T) {
	t.Run("join and leave", func(t *testing.T) {
		room := newRoom()
		go room.Run()
		defer room.Close()

		users := []userMock{newUser(), newUser(), newUser()}
		for _, user := range users {
			room.join <- user
		}
		if len(room.users) != len(users) {
			t.Errorf("join failed: got %v want %v", len(room.users), len(users))
		}

		for _, user := range users {
			room.leave <- user
		}
		if len(room.users) != 0 {
			t.Errorf("leave failed: got %v want 0", len(room.users))
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
