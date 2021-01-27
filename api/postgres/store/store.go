package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required
	"github.com/syncedvideo/syncedvideo"
)

// Store implements syncedvideo.Store
type Store struct {
	user     *UserStore
	room     *RoomStore
	playlist *PlaylistStore
}

// NewStore returns a new store
func NewStore(dataSourceName string) (syncedvideo.Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	userStore := &UserStore{db}
	playlistStore := &PlaylistStore{db}
	roomStore := &RoomStore{db: db, playlist: playlistStore}

	return &Store{
		user:     userStore,
		room:     roomStore,
		playlist: playlistStore,
	}, nil
}

func (s *Store) User() syncedvideo.UserStore {
	return s.user
}

func (s *Store) Room() syncedvideo.RoomStore {
	return s.room
}

func (s *Store) Playlist() syncedvideo.PlaylistStore {
	return s.playlist
}
