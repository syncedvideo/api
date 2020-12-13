package room

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ConnectionHub manages websocket connections
type ConnectionHub struct {
	Connections   map[uuid.UUID]*Connection `json:"connections"`
	ConnectionCap int                       `json:"connectionHub"`
}

// Connection represents a connected user
type Connection struct {
	User          *User             `json:"user"`
	WsConnections []*websocket.Conn `json:"-"`
}

// NewConnectionHub returns a new ConnectionHub
func NewConnectionHub(connectionCap int) *ConnectionHub {
	return &ConnectionHub{
		Connections:   make(map[uuid.UUID]*Connection),
		ConnectionCap: connectionCap,
	}
}

// ErrConnectionCapReached handle connection cap error
var ErrConnectionCapReached = errors.New("Reached connection cap")

// Connect a user
func (hub *ConnectionHub) Connect(user *User, wsConn *websocket.Conn) (*Connection, error) {
	if hub.ConnectionCap == len(hub.Connections) {
		return nil, ErrConnectionCapReached
	}

	connection, exists := hub.Connections[user.ID]
	if exists {
		connection.WsConnections = append(hub.Connections[user.ID].WsConnections, wsConn)
	} else {
		hub.Connections[user.ID] = &Connection{
			User:          user,
			WsConnections: []*websocket.Conn{wsConn},
		}
	}

	log.Printf("User %v connected\n", user.ID)
	return hub.Connections[user.ID], nil
}

// Disconnect a user
func (hub *ConnectionHub) Disconnect(user *User, wsConn *websocket.Conn) {
	connection, exists := hub.Connections[user.ID]
	if !exists {
		return
	}

	// Remove connection from hub if all connections are closed
	if len(connection.WsConnections) == 1 {
		delete(hub.Connections, user.ID)
		log.Printf("User %v disconnected\n", user.ID)
		return
	}

	// Remove WebSocket connection
	for i, v := range connection.WsConnections {
		if v == wsConn {
			connection.WsConnections = append(connection.WsConnections[:i], connection.WsConnections[i+1:]...)
		}
	}
}

// GetAllWebSocketConnections returns all connected websocket connections
func (hub *ConnectionHub) GetAllWebSocketConnections() []*websocket.Conn {
	wsConns := []*websocket.Conn{}
	for _, conn := range hub.Connections {
		wsConns = append(wsConns, conn.WsConnections...)
	}
	return wsConns
}

// BroadcastEvent to all all user connections
func (hub *ConnectionHub) BroadcastEvent(event WsEvent) {
	for _, conn := range hub.Connections {
		for _, wsConn := range conn.WsConnections {
			event.User = conn.User
			eventMsg := new(bytes.Buffer)
			err := json.NewEncoder(eventMsg).Encode(event)
			if err != nil {
				log.Println("BroadcastEvent error:", err)
				continue
			}
			wsConn.WriteMessage(websocket.TextMessage, eventMsg.Bytes())
		}
	}
}

// WsEventName represents name of event
type WsEventName string

// WebSocket events
const (
	WsEventJoin  = WsEventName("join")
	WsEventLeave = WsEventName("leave")
	WsEventSync  = WsEventName("sync")
)

// WsEvent is sent to users
type WsEvent struct {
	Name WsEventName `json:"event"`
	User *User       `json:"user"`
	Data interface{} `json:"data"`
}

// WsActionName represents name of action
type WsActionName string

// WebSocket actions
const (
	// User actions
	WsActionUserSetUsername = WsActionName("user:set:username")
	WsActionUserSetColor    = WsActionName("user:set:color")

	// Player actions
	WsActionPlayerInit          = WsActionName("player:init")
	WsActionPlayerTogglePlaying = WsActionName("player:togglePlaying")
	WsActionPlayerSkip          = WsActionName("player:skip")

	// Queue actions
	WsActionQueueAdd    = WsActionName("queue:add")
	WsActionQueueRemove = WsActionName("queue:remove")
	WsActionQueueVote   = WsActionName("queue:vote")

	// Chat actions
	WsActionChatMessage = WsActionName("chat:message")
)

// WsAction is sent by user
type WsAction struct {
	Name WsActionName    `json:"action"`
	Data json.RawMessage `json:"data"`
	User *User           `json:"-"`
}
