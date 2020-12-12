package room

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ConnectionHub manages websocket connections
type ConnectionHub struct {
	Connections map[uuid.UUID]*Connection
}

// Connection represents a connected user
type Connection struct {
	User          *User
	WsConnections []*websocket.Conn
}

// NewConnectionHub returns a new ConnectionHub
func NewConnectionHub() *ConnectionHub {
	return &ConnectionHub{
		Connections: make(map[uuid.UUID]*Connection),
	}
}

// Connect a user
func (hub *ConnectionHub) Connect(user *User, wsConn *websocket.Conn) *Connection {
	connection, exists := hub.Connections[user.ID]
	if exists {
		connection.WsConnections = append(hub.Connections[user.ID].WsConnections, wsConn)
		log.Println("Connected user %w to ConnectionHub", user.ID)
		return connection
	}
	hub.Connections[user.ID] = &Connection{
		User:          user,
		WsConnections: []*websocket.Conn{wsConn},
	}
	return hub.Connections[user.ID]
}

// Disconnect a user
func (hub *ConnectionHub) Disconnect(user *User) {
	connection, exists := hub.Connections[user.ID]
	if !exists {
		return
	}

	// Remove connection from hub if all connections are closed
	if len(connection.WsConnections) == 0 {
		delete(hub.Connections, user.ID)
		log.Println("Disconnected user %w from ConnectionHub", user.ID)
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
