package syncedvideo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/go-redis/redis/v8"
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
	return room, nil
}

type RoomEvent struct {
	ID string          `json:"id"`
	T  int             `json:"t"`
	D  json.RawMessage `json:"d"`
}

func NewRoomEvent(eventType int, data []byte) RoomEvent {
	return RoomEvent{
		ID: uuid.NewString(),
		T:  eventType,
		D:  data,
	}
}

func (e *RoomEvent) ResetIDFields() {
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

type RedisRoomPubSub struct {
	client *redis.Client
}

func NewRedisRoomPubSub(client *redis.Client) *RedisRoomPubSub {
	pubSub := new(RedisRoomPubSub)
	pubSub.client = client
	return pubSub
}

func (r *RedisRoomPubSub) Publish(roomID string, event RoomEvent) {
	log.Printf("publish %v", event)

	eb, err := json.Marshal(event)
	if err != nil {
		log.Printf("error marshalling event bytes: %v\n", err)
		return
	}

	r.client.Publish(context.Background(), roomID, eb)
}

func (r *RedisRoomPubSub) Subscribe(roomID string) <-chan RoomEvent {
	pubSub := r.client.Subscribe(context.Background(), roomID)
	ch := make(chan RoomEvent)

	go func() {
		for msg := range pubSub.Channel() {
			log.Printf("redis received message: %v\n", msg)

			event := RoomEvent{}
			err := json.Unmarshal([]byte(msg.Payload), &event)
			if err != nil {
				log.Printf("error unmarshalling event: %v\n", err)
				break
			}

			ch <- event
		}
	}()

	return ch
}
