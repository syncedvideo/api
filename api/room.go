package syncedvideo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	Users []User `json:"users"`
	// Users map[uuid.UUID]*User `json:"users"`
	// broadcast  chan []byte
	// register   chan *User
	// unregister chan *User

	store Store
	redis *redis.Client
}

func (r *Room) Run(user *User, store Store, redis *redis.Client) {
	pubsub := redis.Subscribe(context.Background(), r.ID.String())
	defer func() {
		user.conn.Close()
		pubsub.Close()
		err := store.Room().Leave(r, user)
		if err != nil {
			log.Printf("error leaving room: %v", err)
		}
	}()

	err := store.Room().Join(r, user)
	if err != nil {
		log.Printf("error joining room: %v", err)
		return
	}

	// handle incoming messages
	go func() {
		for msg := range pubsub.Channel() {
			fmt.Printf("received message: %s\n", msg)
			user.conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		}
	}()

	// ping to keep connection alive
	go func() {
		for {
			time.Sleep(time.Second * 5)
			err := user.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
			if err != nil {
				break
			}
		}
	}()

	// ws connection loop
	for {
		_, msg, err := user.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			log.Printf("error reading message: %v\n", err)
			break
		}
		log.Printf("recieved message: %v\n", msg)
	}
}

const (
	MessageRoomJoin  = 1000
	MessageRoomLeave = 1001
	MessageRoomChat  = 2000
)

type Message struct {
	T int         `json:"t"`
	D interface{} `json:"d"`
}

func (r *Room) Publish(redis *redis.Client, msgType int, msgData interface{}) error {
	msg := Message{
		T: msgType,
		D: msgData,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	redis.Publish(context.Background(), r.ID.String(), b)
	log.Printf("published: %v\n", msg)
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
