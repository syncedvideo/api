package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/syncedvideo/syncedvideo"
)

// RoomStore implements syncedvideo.RoomStore
type RoomStore struct {
	db *sqlx.DB
}

func (rs *RoomStore) Get(id uuid.UUID) (syncedvideo.Room, error) {
	r := syncedvideo.Room{}
	err := rs.db.Get(&r, `SELECT * FROM sv_room WHERE id=$1 RETURNING *`, id)
	if err != nil {
		return syncedvideo.Room{}, fmt.Errorf("error getting room: %w", err)
	}
	return r, nil
}

func (rs *RoomStore) Create(r *syncedvideo.Room) error {
	err := rs.db.Get(r, `INSERT INTO sv_room VALUES ($1, $2, $3) RETURNING *`, r.ID, nil, r.Name)
	if err != nil {
		return fmt.Errorf("error creating room: %w", err)
	}
	return nil
}

func (rs *RoomStore) Update(r *syncedvideo.Room) error {
	err := rs.db.Get(r, `UPDATE sv_room SET name=$1, owner_user_id=$2 WHERE id=$3 RETURNING *`, r.Name, r.OwnerUserID, r.ID)
	if err != nil {
		return fmt.Errorf("error updating room: %w", err)
	}
	return nil
}

func (rs *RoomStore) Delete(id uuid.UUID) error {
	_, err := rs.db.Exec(`DELETE FROM sv_room WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("error deleting room: %w", err)
	}
	return nil
}

func (rs *RoomStore) GetPlaylistItem(r *syncedvideo.Room, id uuid.UUID) (syncedvideo.PlaylistItem, error) {
	item := syncedvideo.PlaylistItem{}
	err := rs.db.Get(&item, `SELECT * FROM sv_room_playlist_item WHERE room_id=$1 AND playlist_item_id=$2`, r.ID, id)
	if err != nil {
		return syncedvideo.PlaylistItem{}, fmt.Errorf("error getting playlist item: %w", err)
	}
	return item, nil
}

func (rs *RoomStore) GetAllPlaylistItems(r *syncedvideo.Room) ([]syncedvideo.PlaylistItem, error) {
	items := []syncedvideo.PlaylistItem{}
	err := rs.db.Select(&items, `SELECT * FROM sv_room_playlist_item WHERE room_id=$1`, r.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting all playlist items: %w", err)
	}
	return items, nil
}

func (rs *RoomStore) CreatePlaylistItem(r *syncedvideo.Room, p *syncedvideo.PlaylistItem) error {
	err := rs.db.Get(p, `INSERT INTO sv_room_playlist_item VALUES ($1, $2, $3) RETURNING *`, p.ID, r.ID, p.UserID)
	if err != nil {
		return fmt.Errorf("error creating playlist item: %w", err)
	}
	return nil
}

func (rs *RoomStore) UpdatePlaylistItem(r *syncedvideo.Room, p *syncedvideo.PlaylistItem) error {
	err := rs.db.Get(p, `UPDATE sv_room_playlist_item SET room_id=$1, user_id=$2 WHERE id=$3`, r.ID, p.UserID, p.ID)
	if err != nil {
		return fmt.Errorf("error updating playlist item: %w", err)
	}
	return nil
}

func (rs *RoomStore) DeletePlaylistItem(r *syncedvideo.Room, id uuid.UUID) error {
	_, err := rs.db.Exec(`DELETE FROM sv_room_playlist_item WHERE id=$1 AND room_id=$2`, id, r.ID)
	if err != nil {
		return fmt.Errorf("error deleting playlist item: %w", err)
	}
	return nil
}
