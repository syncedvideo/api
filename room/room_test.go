package room

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestRun(t *testing.T) {
	t.Run("register and unregister", func(t *testing.T) {
		room := newRoom()
		go room.Run()
		defer room.Close()

		users := []userMock{newUser(), newUser(), newUser()}
		for _, user := range users {
			room.register <- user
		}
		if len(room.users) != len(users) {
			t.Errorf("register failed: got %v want %v", len(room.users), len(users))
		}

		for _, user := range users {
			room.unregister <- user
		}
		if len(room.users) != 0 {
			t.Errorf("unregister failed: got %v want 0", len(room.users))
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
