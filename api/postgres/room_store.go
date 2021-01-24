package postgres

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/syncedvideo/syncedvideo"
)

// RoomStore implements syncedvideo.RoomStore
type RoomStore struct {
	db *sqlx.DB
}

func (rs *RoomStore) Get(id uuid.UUID) (syncedvideo.Room, error) {
	var r syncedvideo.Room
	err := rs.db.Get(&r, `SELECT * FROM sv_room WHERE id = $1 RETURNING *`, id)
	if err != nil {
		return syncedvideo.Room{}, err
	}
	return r, nil
}

func (rs *RoomStore) Create(r *syncedvideo.Room) error {
	err := rs.db.Get(r, `INSERT INTO sv_room VALUES ($1, $2, $3) RETURNING *`, r.ID, nil, r.Name)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoomStore) Update(r *syncedvideo.Room) error {
	err := rs.db.Get(r, `UPDATE sv_room SET name = $1, owner_user_id = $2 WHERE id = $3 RETURNING *`, r.Name, r.OwnerUserID, r.ID)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoomStore) Delete(id uuid.UUID) error {
	_, err := rs.db.Exec(`DELETE FROM sv_room WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
