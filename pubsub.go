package syncedvideo

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type PubSub interface {
	Publish(roomID string, event Event)
	Subscribe(roomID string) <-chan Event
}

type RedisPubSub struct {
	client *redis.Client
}

func NewRedisPubSub(client *redis.Client) *RedisPubSub {
	pubSub := new(RedisPubSub)
	pubSub.client = client
	return pubSub
}

func (r *RedisPubSub) Publish(roomID string, event Event) {
	log.Printf("publish %v", event)

	eb, err := json.Marshal(event)
	if err != nil {
		log.Printf("error marshalling event bytes: %v\n", err)
		return
	}

	r.client.Publish(context.Background(), roomID, eb)
}

func (r *RedisPubSub) Subscribe(roomID string) <-chan Event {
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
