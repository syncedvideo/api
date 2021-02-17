package postgres

import (
	"testing"

	"github.com/syncedvideo/syncedvideo"
)

func TestRoomStore(t *testing.T) {
	testUser := syncedvideo.User{}
	createUser(&testUser)

	testRoom := syncedvideo.Room{OwnerUserID: testUser.ID}
	err := createRoom(&testRoom)
	if err != nil {
		t.Errorf("create room failed: %w", err)
	}

	t.Run("get room", func(t *testing.T) {
		room, err := store.Room().Get(testRoom.ID)
		if err != nil {
			t.Errorf("Get failed: %w", err)
		}
		if room.ID != testRoom.ID {
			t.Errorf("got %s want %s", room.ID, testRoom.ID)
		}
	})

	t.Run("update room", func(t *testing.T) {
		room := syncedvideo.Room{ID: testRoom.ID}
		room.Name = "Test"
		err := store.Room().Update(&room)
		if err != nil {
			t.Errorf("Update failed: %w", err)
		}
	})

	t.Run("join room", func(t *testing.T) {
		err := store.Room().Join(&testRoom, &testUser)
		if err != nil {
			t.Errorf("Join failed: %w", err)
		}
		users, err := store.Room().GetUsers(&testRoom)
		if err != nil {
			t.Errorf("GetUsers failed: %w", err)
		}
		if len(users) != 1 {
			t.Errorf("got %v want 1", len(users))
		}
	})

	t.Run("leave room", func(t *testing.T) {
		err := store.Room().Leave(&testRoom, &testUser)
		if err != nil {
			t.Errorf("Leave failed: %w", err)
		}
		users, err := store.Room().GetUsers(&testRoom)
		if err != nil {
			t.Errorf("GetUsers failed: %w", err)
		}
		if len(users) != 0 {
			t.Errorf("got %v want 0", len(users))
		}
	})

	t.Run("delete room", func(t *testing.T) {
		err := store.Room().Delete(testRoom.ID)
		if err != nil {
			t.Errorf("Delete failed: %w", err)
		}
	})

	deleteUser(testUser.ID)
}

func createRoom(room *syncedvideo.Room) error {
	return store.Room().Create(room)
}
