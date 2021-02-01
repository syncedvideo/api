package syncedvideo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID            uuid.UUID                  `db:"id" json:"id"`
	Name          string                     `db:"name" json:"name"`
	Description   string                     `db:"description" json:"description"`
	OwnerUserID   uuid.UUID                  `json:"ownerUserId" db:"owner_user_id"`
	PlaylistItems map[uuid.UUID]PlaylistItem `json:"playlistItems"`

	Users map[uuid.UUID]*User `json:"users"`
	// broadcast  chan []byte
	// register   chan *User
	// unregister chan *User

	store Store
	redis *redis.Client
}

func (r *Room) Run(user *User, store Store, redis *redis.Client) {
	r.store = store
	r.redis = redis
	pubsub := r.redis.Subscribe(context.Background(), r.ID.String())

	go func() {
		for {
			_, msg, err := user.conn.ReadMessage()
			if err != nil {
				log.Printf("error reading message: %v\n", err)
				break
			}
			log.Printf("recieved message: %v\n", msg)
		}
	}()

	for msg := range pubsub.Channel() {
		fmt.Printf("received message: %s\n", msg)
		user.conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
	}
}

func (r *Room) Publish(redis *redis.Client, message interface{}) error {
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	redis.Publish(context.Background(), r.ID.String(), b)
	fmt.Printf("published: %s\n", message)
	return nil
}

type PlaylistItem struct {
	ID     uuid.UUID `db:"id"`
	RoomID uuid.UUID `db:"room_id"`
	UserID uuid.UUID `db:"user_id"`
	Votes  []PlaylistItemVote
}

type PlaylistItemVote struct {
	ID     uuid.UUID `db:"id"`
	ItemID uuid.UUID `db:"item_id"`
	UserID uuid.UUID `db:"user_id"`
}

func NewRoom(connectionCap int) *Room {
	return &Room{
		ID: uuid.New(),
	}
}
