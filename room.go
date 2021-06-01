package syncedvideo

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewRoom(reader io.Reader) (Room, error) {
	room := Room{}
	err := json.NewDecoder(reader).Decode(&room)
	if err != nil {
		return room, fmt.Errorf("error decoding room: %v", err)
	}
	room.ID = uuid.NewString()
	return room, nil
}
