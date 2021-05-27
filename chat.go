package syncedvideo

import (
	"github.com/google/uuid"
)

type ChatMessage struct {
	ID      string `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

func NewChatMessage(author, message string) ChatMessage {
	return ChatMessage{
		ID:      uuid.NewString(),
		Author:  author,
		Message: message,
	}
}
