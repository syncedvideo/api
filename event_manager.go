package syncedvideo

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type EventManager interface {
	Publish(roomID string, event Event)
	Subscribe(roomID string) <-chan Event
}

type Event struct {
	T EventType       `json:"t"`
	D json.RawMessage `json:"d"`
}

type EventType string

func (e EventType) String() string {
	return string(e)
}

var (
	EventChat EventType = "chat"
	EventPlay EventType = "play"
)

func NewEvent(eventType EventType, data interface{}) Event {
	dataB, _ := json.Marshal(data)
	return Event{
		T: eventType,
		D: dataB,
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
