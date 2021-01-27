package store

import (
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
	panic("not implemented") // TODO: Implement
}

func (s *PlaylistStore) Create(p *syncedvideo.PlaylistItem) error {
	panic("not implemented") // TODO: Implement
}

func (s *PlaylistStore) Update(p *syncedvideo.PlaylistItem) error {
	panic("not implemented") // TODO: Implement
}

func (s *PlaylistStore) Delete(id uuid.UUID) error {
	panic("not implemented") // TODO: Implement
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
