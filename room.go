package syncedvideo

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewRoom(reader io.Reader) (Room, error) {
	room := Room{}
	err := json.NewDecoder(reader).Decode(&room)
	if err != nil {
		return room, fmt.Errorf("error decoding room: %v", err)
	}
	return room, nil
}
