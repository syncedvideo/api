package postgres

import (
	"testing"

	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
)

func TestUserStore(t *testing.T) {
	userID := uuid.New()
	t.Run("create user", func(t *testing.T) {
		user := syncedvideo.User{ID: userID}
		err := createUser(&user)
		if err != nil {
			t.Errorf("Create failed: %w", err)
		}
	})
	t.Run("get user", func(t *testing.T) {
		_, err := store.User().Get(userID)
		if err != nil {
			t.Errorf("Get failed: %w", err)
		}
	})
	t.Run("update user", func(t *testing.T) {
		user := syncedvideo.User{ID: userID}
		user.Name = "TestUser"
		err := store.User().Update(&user)
		if err != nil {
			t.Errorf("Update failed: %w", err)
		}
	})
	t.Run("delete user", func(t *testing.T) {
		err := deleteUser(userID)
		if err != nil {
			t.Errorf("Delete failed: %w", err)
		}
	})
}

func createUser(user *syncedvideo.User) error {
	return store.User().Create(user)
}

func deleteUser(userID uuid.UUID) error {
	return store.User().Delete(userID)
}
