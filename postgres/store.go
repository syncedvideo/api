package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // required
)

// Store manages all DB connections
type Store struct {
	//
}

// NewStore returns a new store
func NewStore(dataSourceName string) (*Store, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return &Store{}, nil
}
