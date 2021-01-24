package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required
)

// Store implements syncedvideo.Store
type Store struct {
	*RoomStore
	*UserStore
}

// NewStore returns a new store
func NewStore(dataSourceName string) (Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return Store{}, fmt.Errorf("error opening database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return Store{}, fmt.Errorf("error connecting to database: %w", err)
	}
	return Store{
		&RoomStore{db},
		&UserStore{db},
	}, nil
}
