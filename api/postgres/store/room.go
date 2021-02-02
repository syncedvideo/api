package store

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/syncedvideo/syncedvideo"
)

// RoomStore implements syncedvideo.RoomStore
type RoomStore struct {
	db       *sqlx.DB
	playlist *PlaylistStore
}

func (s *RoomStore) Get(roomID uuid.UUID) (syncedvideo.Room, error) {
	r := syncedvideo.Room{}
	err := s.db.Get(&r, "SELECT * FROM sv_room AS room WHERE room.id = $1", roomID)
	if err == sql.ErrNoRows {
		return syncedvideo.Room{}, err
	} else if err != nil {
		return syncedvideo.Room{}, fmt.Errorf("error getting room: %w", err)
	}
	items, err := s.playlist.All(r.ID)
	if err != nil {
		return syncedvideo.Room{}, fmt.Errorf("error getting room playlist items: %s", err)
	}
	r.PlaylistItems = items
	return r, nil
}

func (s *RoomStore) Create(r *syncedvideo.Room) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.OwnerUserID == uuid.Nil {
		return errors.New("OwnerUserID is required")
	}
	err := s.db.Get(r, "INSERT INTO sv_room VALUES ($1, $2, $3, $4) RETURNING *", r.ID, r.OwnerUserID, r.Name, r.Description)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return fmt.Errorf("error creating room: %w", err)
	}
	return nil
}

func (s *RoomStore) Update(r *syncedvideo.Room) error {
	err := s.db.Get(r, "UPDATE sv_room SET name=$1, owner_user_id=$2 WHERE id=$3 RETURNING *", r.Name, r.OwnerUserID, r.ID)
	if err != nil {
		return fmt.Errorf("error updating room: %w", err)
	}
	return nil
}

func (s *RoomStore) Delete(roomID uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM sv_room WHERE id=$1", roomID)
	if err != nil {
		return fmt.Errorf("error deleting room: %w", err)
	}
	return nil
}

func (s *RoomStore) GetPlaylistItem(r *syncedvideo.Room, id uuid.UUID) (syncedvideo.PlaylistItem, error) {
	item := syncedvideo.PlaylistItem{}
	err := s.db.Get(&item, "SELECT * FROM sv_playlist_item WHERE room_id=$1 AND item_id=$2", r.ID, id)
	if err != nil {
		return syncedvideo.PlaylistItem{}, fmt.Errorf("error getting playlist item: %w", err)
	}
	return item, nil
}

func (s *RoomStore) Join(r *syncedvideo.Room, u *syncedvideo.User) error {
	ru := struct {
		RoomID uuid.UUID `db:"room_id"`
		UserID uuid.UUID `db:"user_id"`
	}{}
	err := s.db.Get(&ru, "SELECT * from sv_room_user WHERE room_id = $1 AND user_id = $2", r.ID, u.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if ru.RoomID != uuid.Nil || ru.UserID != uuid.Nil {
		err = s.Leave(r, u)
		if err != nil {
			return err
		}
	}

	_, err = s.db.Exec("INSERT INTO sv_room_user VALUES($1, $2)", r.ID, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoomStore) Leave(r *syncedvideo.Room, u *syncedvideo.User) error {
	_, err := s.db.Exec("DELETE FROM sv_room_user WHERE room_id = $1 AND user_id = $2", r.ID, u.ID)
	if err != nil {
		return err
	}
	return nil
}
