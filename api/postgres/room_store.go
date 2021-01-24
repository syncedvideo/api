package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/syncedvideo/syncedvideo"
)

// RoomStore implements syncedvideo.RoomStore
type RoomStore struct {
	*sqlx.DB
}

func (db *RoomStore) GetRoom(id uuid.UUID) (syncedvideo.Room, error) {
	var r syncedvideo.Room
	err := db.Get(&r, `SELECT * FROM sv_room WHERE id = $1 RETURNING *`, id)
	if err != nil {
		return syncedvideo.Room{}, err
	}
	return r, nil
}

func (db RoomStore) CreateRoom(r *syncedvideo.Room) error {
	err := db.Get(r, `INSERT INTO sv_room VALUES ($1, $2, $3) RETURNING *`, r.ID, nil, r.Name)
	if err != nil {
		return err
	}
	return nil
}

func (db *RoomStore) UpdateRoom(r *syncedvideo.Room) error {
	err := db.Get(r, `UPDATE sv_room SET name = $1, owner_user_id = $2 WHERE id = $3 RETURNING *`, r.Name, r.OwnerUserID, r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *RoomStore) DeleteRoom(id uuid.UUID) error {
	_, err := db.Exec(`DELETE FROM sv_room WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
