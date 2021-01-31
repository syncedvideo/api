package syncedvideo

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	Messages []*ChatMessage `json:"messages"`
}

type ChatMessage struct {
	ID        uuid.UUID `json:"id"`
	User      User      `json:"user"`
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}

func NewChatMessage(user User, text string) *ChatMessage {
	return &ChatMessage{
		ID:        uuid.New(),
		User:      user,
		Timestamp: time.Now(),
		Text:      text,
	}
}

func NewChat() *Chat {
	return &Chat{
		Messages: make([]*ChatMessage, 0),
	}
}

func (c *Chat) NewMessage(user User, text string) *ChatMessage {
	message := &ChatMessage{
		ID:        uuid.New(),
		User:      user,
		Text:      text,
		Timestamp: time.Now(),
	}
	c.Messages = append(c.Messages, message)
	return message
}
