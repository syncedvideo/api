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

func (r *Room) Close(user *User) {
	fmt.Println("defer")
	user.Connection.Close()
	Config.Store.Room().Leave(r, user)
	r.Publish(WebSocketMessageLeave, user)
	r.SyncUsers()
}

func (r *Room) Run(user *User) {
	defer r.Close(user)

	Config.Store.Room().Join(r, user)
	r.Publish(WebSocketMessageJoin, user)
	pubsub := Config.Redis.Subscribe(context.Background(), r.ID.String())

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
			ping, _ := NewWebSocketMessage(WebSocketMessagePing, "ping").MarshalBinary()
			err := user.Connection.WriteMessage(websocket.TextMessage, ping)
			if err != nil {
				break
			}
		}
	}()

	r.SyncUsers()

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
	users, err := Config.Store.Room().GetUsers(r)
	if err != nil {
		return fmt.Errorf("GetUsers failed: %w", err)
	}
	err = r.Publish(WebSocketMessageSyncUsers, users)
	if err != nil {
		return fmt.Errorf("Publish WebSocketMessageSyncUsers failed: %w", err)
	}
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
