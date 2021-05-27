package syncedvideo

import (
	"encoding/json"
	"fmt"
	"io"
)

type ChatMessage struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

func NewChatMessage(reader io.Reader) (ChatMessage, error) {
	message := ChatMessage{}
	err := json.NewDecoder(reader).Decode(&message)
	if err != nil {
		return message, fmt.Errorf("error decoding chat message: %v", err)
	}
	return message, nil
}
