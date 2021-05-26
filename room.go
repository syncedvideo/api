package syncedvideo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/go-redis/redis/v8"
)

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Chat *Chat  `json:"chat"`
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
	T int             `json:"t"`
	D json.RawMessage `json:"d"`
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
			log.Printf("received %v", msg)

			event := RoomEvent{}
			err := json.Unmarshal([]byte(msg.Payload), &event)
			if err != nil {
				log.Printf("error unmarshalling event: %v", err)
				break
			}

			ch <- event
		}
	}()

	return ch
}

type Chat struct {
	Messages []ChatMessage `json:"messages"`
}

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
