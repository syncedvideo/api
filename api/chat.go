package syncedvideo

import (
	"time"

	"github.com/google/uuid"
)

// Chat represents the room's chat
type Chat struct {
	Messages []*ChatMessage `json:"messages"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        uuid.UUID `json:"id"`
	User      *User     `json:"user"`
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}

// NewChat returns a new chat
func NewChat() *Chat {
	return &Chat{
		Messages: make([]*ChatMessage, 0),
	}
}

// NewMessage adds a new chat message
func (c *Chat) NewMessage(user *User, text string) *ChatMessage {
	message := &ChatMessage{
		ID:        uuid.New(),
		User:      user,
		Text:      text,
		Timestamp: time.Now(),
	}
	c.Messages = append(c.Messages, message)
	return message
}
