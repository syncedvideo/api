package syncedvideo

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	WebSocketMessagePing      = 0
	WebSocketMessageJoin      = 1000
	WebSocketMessageLeave     = 1001
	WebSocketMessageChat      = 2000
	WebSocketMessageSyncUsers = 3000
)

type WebSocketMessage struct {
	ID uuid.UUID   `json:"id"`
	T  int         `json:"t"`
	D  interface{} `json:"d"`
}

func NewWebSocketMessage(msgType int, msgData interface{}) *WebSocketMessage {
	return &WebSocketMessage{
		ID: uuid.New(),
		T:  msgType,
		D:  msgData,
	}
}

// MarshalBinary implements encoding.BinaryMarshaler
func (msg *WebSocketMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}
