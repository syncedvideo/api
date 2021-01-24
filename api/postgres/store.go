package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required
)

// Store manages all DB connections
type Store struct {
	*UserStore
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
		&UserStore{db},
	}, nil
}
