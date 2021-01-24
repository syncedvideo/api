package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/syncedvideo/syncedvideo"
)

type UserStore struct {
	*sqlx.DB
}

func (db *UserStore) GetUser(id uuid.UUID) (syncedvideo.User, error) {
	u := syncedvideo.User{}
	err := db.Get(&u, `SELECT * FROM sv_user where id = $1`, id)
	if err != nil {
		return syncedvideo.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (db *UserStore) CreateUser(u *syncedvideo.User) error {
	createdAt := time.Now().UTC()
	err := db.Get(u, `
		INSERT INTO sv_user 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`, u.ID, u.Name, u.Color, u.IsAdmin, createdAt)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (db *UserStore) UpdateUser(u *syncedvideo.User) error {
	err := db.Get(u, `
		UPDATE sv_user
		SET name = $1, color = $2, is_admin = $3, updated_at = $4
		WHERE id = $5
		RETURNING *
	`, u.Name, u.Color, u.IsAdmin, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (db *UserStore) DeleteUser(id uuid.UUID) error {
	_, err := db.Exec(`DELETE from sv_user WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
