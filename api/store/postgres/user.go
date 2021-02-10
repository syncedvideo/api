package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/syncedvideo/syncedvideo"
)

type UserStore struct {
	db *sqlx.DB
}

func (s *UserStore) Get(userID uuid.UUID) (syncedvideo.User, error) {
	u := syncedvideo.User{}
	err := s.db.Get(&u, `SELECT * FROM sv_user where id=$1`, userID)
	if err != nil {
		return syncedvideo.User{}, err
	}
	return u, nil
}

func (s *UserStore) Create(u *syncedvideo.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	err := s.db.Get(u, `INSERT INTO sv_user VALUES ($1, $2, $3, $4) RETURNING *`, u.ID, u.Name, u.Color, u.IsAdmin)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (s *UserStore) Update(u *syncedvideo.User) error {
	err := s.db.Get(u, `UPDATE sv_user SET name=$1, color=$2, is_admin=$3 WHERE id=$4 RETURNING *`, u.Name, u.Color, u.IsAdmin, u.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) Delete(userID uuid.UUID) error {
	_, err := s.db.Exec(`DELETE from sv_user WHERE id=$1`, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (s *UserStore) GetCurrentRooms(userID uuid.UUID) ([]syncedvideo.Room, error) {
	rooms := []syncedvideo.Room{}
	err := s.db.Select(&rooms, `
		SELECT sr.*
		FROM sv_room_user_connection sruc 
		LEFT JOIN sv_room sr ON sruc.room_id = sr.id 
		WHERE user_id = $1
	`, userID)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, fmt.Errorf("error getting current rooms: %w", err)
	}
	return rooms, nil
}
