package store

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/syncedvideo/syncedvideo"
)

// PlaylistStore implements syncedvideo.PlaylistStore
type PlaylistStore struct {
	db *sqlx.DB
}

func (s *PlaylistStore) Get(id uuid.UUID) (syncedvideo.PlaylistItem, error) {
	panic("not implemented") // TODO: Implement
}

func (s *PlaylistStore) All(roomID uuid.UUID) (map[uuid.UUID]syncedvideo.PlaylistItem, error) {
	rows, err := s.db.Query(`
		SELECT item.id, item.room_id, item.user_id, vote.id, vote.item_id, vote.user_id 
		FROM sv_playlist_item item
		LEFT JOIN sv_playlist_item_vote vote
		ON item.id = vote.item_id
		WHERE item.room_id = $1
	`, roomID)
	if err != nil {
		return nil, fmt.Errorf("error getting all playlist items: %w", err)
	}

	items := make(map[uuid.UUID]syncedvideo.PlaylistItem)
	for rows.Next() {
		item := syncedvideo.PlaylistItem{}
		vote := syncedvideo.PlaylistItemVote{}
		err := rows.Scan(&item.ID, &item.RoomID, &item.UserID, &vote.ID, &vote.ItemID, &vote.UserID)
		if err != nil {
			return nil, fmt.Errorf("error scanning room playlist item: %w", err)
		}
		if vote.ID != uuid.Nil {
			item.Votes = append(item.Votes, vote)
		}
		if _, exists := items[item.ID]; !exists {
			items[item.ID] = item
		}
	}

	return items, nil
}

func (s *PlaylistStore) Create(p *syncedvideo.PlaylistItem) error {
	err := s.db.Get(p, "INSERT INTO sv_playlist_item VALUES ($1, $2, $3) RETURNING *", p.ID, p.RoomID, p.UserID)
	if err != nil {
		return fmt.Errorf("error creating playlist item: %w", err)
	}
	return nil
}

func (s *PlaylistStore) Update(p *syncedvideo.PlaylistItem) error {
	err := s.db.Get(p, "UPDATE sv_playlist_item SET room_id=$1, user_id=$2 WHERE id=$3 RETURNING *", p.RoomID, p.UserID, p.ID)
	if err != nil {
		return fmt.Errorf("error updating playlist item: %w", err)
	}
	return nil
}

func (s *PlaylistStore) Delete(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM sv_playlist_item WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error deleting playlist item: %w", err)
	}
	return nil
}

func (s *PlaylistStore) GetVote(itemID uuid.UUID, userID uuid.UUID) (syncedvideo.PlaylistItemVote, error) {
	panic("not implemented") // TODO: Implement
}

func (s *PlaylistStore) CreateVote(itemID uuid.UUID, userID uuid.UUID) error {
	panic("not implemented") // TODO: Implement
}

func (s *PlaylistStore) DeleteVote(itemID uuid.UUID, userID uuid.UUID) error {
	panic("not implemented") // TODO: Implement
}
