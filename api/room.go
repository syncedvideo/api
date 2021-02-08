package syncedvideo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	ID            uuid.UUID                  `db:"id" json:"id"`
	Name          string                     `db:"name" json:"name"`
	Description   string                     `db:"description" json:"description"`
	OwnerUserID   uuid.UUID                  `json:"ownerUserId" db:"owner_user_id"`
	PlaylistItems map[uuid.UUID]PlaylistItem `json:"playlistItems"`

	Users []User `json:"users"`
	// Users map[uuid.UUID]*User `json:"users"`
	// broadcast  chan []byte
	// register   chan *User
	// unregister chan *User
}

func (r *Room) Run(user *User) {
	pubsub := Config.Redis.Subscribe(context.Background(), r.ID.String())
	defer func() {
		pubsub.Close()
	}()
	go func() {
		for msg := range pubsub.Channel() {
			fmt.Printf("received message: %s\n", msg)
			user.Connection.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		}
	}()

	// ping to keep connection alive
	go func() {
		for {
			time.Sleep(time.Second * 5)
			err := user.Connection.WriteMessage(websocket.TextMessage, []byte("ping"))
			if err != nil {
				break
			}
		}
	}()

	// ws connection loop
	for {
		_, msg, err := user.Connection.ReadMessage()
		if err != nil {
			log.Println(err)
			log.Printf("error reading message: %v\n", err)
			break
		}
		log.Printf("recieved message: %v\n", msg)
	}
}

func (r *Room) Publish(msgType int, msgData interface{}) error {
	msg := NewWebSocketMessage(msgType, msgData)
	Config.Redis.Publish(context.Background(), r.ID.String(), msg)
	log.Printf("published: %v\n", msg)
	return nil
}

func (r *Room) SyncUsers() error {
	return r.Publish(WebSocketMessageSyncUsers, r.Users)
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
