package syncedvideo

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Event struct {
	ID string          `json:"id"`
	T  int             `json:"t"`
	D  json.RawMessage `json:"d"`
}

func NewEvent(eventType int, data []byte) Event {
	return Event{
		ID: uuid.NewString(),
		T:  eventType,
		D:  data,
	}
}

func (e *Event) resetIDFields() {
	e.ID = ""

	data := make(map[string]interface{})
	json.Unmarshal(e.D, &data)

	_, ok := data["id"]
	if ok {
		data["id"] = ""
	}

	dataB, _ := json.Marshal(data)
	e.D = dataB
}
