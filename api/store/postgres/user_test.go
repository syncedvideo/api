package postgres

import (
	"testing"

	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
)

func TestUserStore(t *testing.T) {
	id := uuid.New()
	t.Run("create user", func(t *testing.T) {
		user := syncedvideo.User{ID: id}
		err := store.User().Create(&user)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("get user", func(t *testing.T) {
		_, err := store.User().Get(id)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("update user", func(t *testing.T) {
		user := syncedvideo.User{ID: id}
		user.Name = "TestUser"
		err := store.User().Update(&user)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("delete user", func(t *testing.T) {
		err := store.User().Delete(id)
		if err != nil {
			t.Error(err)
		}
	})
}
