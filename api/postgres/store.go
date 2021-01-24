package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required
	"github.com/syncedvideo/syncedvideo"
)

// Store implements syncedvideo.Store
type Store struct {
	user *UserStore
	room *RoomStore
}

// NewStore returns a new store
func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return &Store{
		user: &UserStore{db},
		room: &RoomStore{db},
	}, nil
}

func (s *Store) User() syncedvideo.UserStore {
	return s.user
}

func (s *Store) Room() syncedvideo.RoomStore {
	return s.room
}
