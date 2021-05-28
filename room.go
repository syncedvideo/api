package syncedvideo

import (
	"encoding/json"
	"fmt"
	"io"
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
	return room, nil
}
