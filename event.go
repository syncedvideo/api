package syncedvideo

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type EventManager interface {
	Publish(roomID string, event Event)
	Subscribe(roomID string) <-chan Event
}

type Event struct {
	ID string          `json:"id"`
	T  EventType       `json:"t"`
	D  json.RawMessage `json:"d"`
}

type EventType string

var (
	EventChat EventType = "chat"
)

func NewEvent(eventType EventType, data []byte) Event {
	return Event{
		ID: uuid.NewString(),
		T:  eventType,
		D:  data,
	}
}

type RedisEventManager struct {
	client *redis.Client
}

func NewRedisEventManager(client *redis.Client) *RedisEventManager {
	eventManager := new(RedisEventManager)
	eventManager.client = client
	return eventManager
}

func (r *RedisEventManager) Publish(roomID string, event Event) {
	log.Printf("publish %v", event)

	eventB, err := json.Marshal(event)
	if err != nil {
		log.Printf("error marshalling event: %s\n", err)
		return
	}

	r.client.Publish(context.Background(), roomID, eventB)
}

func (r *RedisEventManager) Subscribe(roomID string) <-chan Event {
	pubSub := r.client.Subscribe(context.Background(), roomID)
	ch := make(chan Event)

	go func() {
		for msg := range pubSub.Channel() {
			log.Printf("redis received message: %v\n", msg)

			event := Event{}
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