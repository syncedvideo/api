package backend

// import (
// 	"bytes"
// 	"encoding/json"
// 	"log"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/gorilla/websocket"
// )

// // Room handles state and connected clients
// type Room struct {
// 	ID      uuid.UUID             `json:"id"`
// 	Clients map[uuid.UUID]*Client `json:"clients"`
// 	Player  *Player               `json:"player"`
// 	Chat    *Chat                 `json:"chat"`
// }

// // Player represents the room's video player
// type Player struct {
// 	CurrentVideo *QueueItem `json:"currentVideo"`
// 	Queue        *Queue     `json:"queue"`
// 	Playing      bool       `json:"playing"`
// }

// // Chat represents the room's chat
// type Chat struct {
// 	Messages []ChatMessage `json:"messages"`
// }

// // ChatMessage represents a chat message
// type ChatMessage struct {
// 	ID        uuid.UUID `json:"id"`
// 	Client    *Client   `json:"client"`
// 	Timestamp time.Time `json:"timestamp"`
// 	Text      string    `json:"text"`
// }

// // Broadcast a RoomEvent to all connected clients
// func (r Room) Broadcast(event RoomEvent) {
// 	for _, client := range r.Clients {
// 		event.Client = *client
// 		eventMsg := new(bytes.Buffer)
// 		err := json.NewEncoder(eventMsg).Encode(event)
// 		if err != nil {
// 			log.Println("Broadcast json.NewEncoder.Encode error:", err)
// 			continue
// 		}
// 		client.Connection.WriteMessage(websocket.TextMessage, eventMsg.Bytes())
// 	}
// }

// // NewRoom creates a new room
// func NewRoom() *Room {
// 	return &Room{
// 		ID:      uuid.New(),
// 		Clients: make(map[uuid.UUID]*Client),
// 		Player: &Player{
// 			CurrentVideo: nil,
// 			Queue: &Queue{
// 				Items: []*QueueItem{},
// 			},
// 		},
// 		Chat: &Chat{},
// 	}
// }

// // Client handles the connection
// type Client struct {
// 	ID         uuid.UUID       `json:"id"`
// 	Connection *websocket.Conn `json:"-"`
// 	Room       Room            `json:"-"`
// 	Username   string          `json:"username"`
// 	Color      string          `json:"color"`
// }

// // Connect to room
// func (c *Client) Connect(conn *websocket.Conn) {
// 	c.Connection = conn
// 	c.Room.Clients[c.ID] = c
// 	c.SyncRoomState()
// }

// // Disconnect from room
// func (c Client) Disconnect() {
// 	delete(c.Room.Clients, c.ID)
// 	c.SyncRoomState()
// }

// // NewClient returns and registers a new client
// func (r Room) NewClient() *Client {
// 	return &Client{
// 		ID:       uuid.New(),
// 		Room:     r,
// 		Username: "",
// 		Color:    "",
// 	}
// }

// // GetConnections returns all open connections
// func (r Room) GetConnections() []*websocket.Conn {
// 	var connections []*websocket.Conn
// 	for _, client := range r.Clients {
// 		log.Println(client.ID)
// 		connections = append(connections, client.Connection)
// 	}
// 	return connections
// }

// // Vote toggles a vote on given QueueItem
// func (c *Client) Vote(queueItemID uuid.UUID) {
// 	// find queue item
// 	foundQueueItem := &QueueItem{}
// 	for _, queueItem := range c.Room.Player.Queue.Items {
// 		if queueItem.ID == queueItemID {
// 			foundQueueItem = queueItem
// 			break
// 		}
// 	}

// 	if foundQueueItem.ID.String() == "" {
// 		log.Println("Client.Vote: Queue item ID not found:", queueItemID)
// 		return
// 	}

// 	_, voted := foundQueueItem.Votes[c.ID]
// 	if !voted {
// 		foundQueueItem.Votes[c.ID] = &Vote{Client: c}
// 		log.Println("Added vote:", c.ID)
// 		return
// 	}
// 	delete(foundQueueItem.Votes, c.ID)
// 	log.Println("Deleted vote:", c.ID)
// }

// // Queue handles the queue of a room
// type Queue struct {
// 	Items []*QueueItem `json:"items"`
// }

// // QueueItem ...
// type QueueItem struct {
// 	ID      uuid.UUID           `json:"id"`
// 	Video   Video               `json:"video"`
// 	Votes   map[uuid.UUID]*Vote `json:"votes"`
// 	AddedBy *Client             `json:"addedBy"`
// }

// // Video represents a video
// type Video struct {
// 	YouTubeID string `json:"youTubeId"`
// 	Title     string `json:"title"`
// }

// // Vote ..
// type Vote struct {
// 	Client *Client `json:"client"`
// }

// // WebSocket events
// const (
// 	RoomEventJoin  = "join"
// 	RoomEventLeave = "leave"
// 	RoomEventSync  = "sync"
// )

// // WebSocket actions
// const (
// 	// Client actions
// 	RoomActionClientSetUsername = "client:set:username"
// 	RoomActionClientSetColor    = "client:set:color"
// 	// Player actions
// 	RoomActionPlayerInit          = "player:init"
// 	RoomActionPlayerTogglePlaying = "player:togglePlaying"
// 	RoomActionPlayerSkip          = "player:skip"
// 	// Chat actions
// 	RoomActionChatMessage = "chat:message"
// 	// Queue actions
// 	RoomActionQueueAdd    = "queue:add"
// 	RoomActionQueueRemove = "queue:remove"
// 	RoomActionQueueVote   = "queue:vote"
// )

// // RoomEvent is sent to client
// type RoomEvent struct {
// 	Event  string      `json:"event"`
// 	Client Client      `json:"client"`
// 	Data   interface{} `json:"data"`
// }

// // RoomAction is sent by client
// type RoomAction struct {
// 	Action string          `json:"action"`
// 	Data   json.RawMessage `json:"data"`
// 	Client *Client         `json:"-"`
// 	Room   *Room           `json:"-"`
// }

// // SyncRoomState broadcasts current room state to all connected clients
// func (c Client) SyncRoomState() {
// 	c.Room.Broadcast(RoomEvent{
// 		Event: RoomEventSync,
// 		Data:  c.Room,
// 	})
// }

// // NewRoomAction returns a new RoomAction
// func NewRoomAction(msg []byte, room *Room, client *Client) (*RoomAction, error) {
// 	ra := RoomAction{
// 		Room:   room,
// 		Client: client,
// 	}
// 	err := json.Unmarshal(msg, &ra)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ra, nil
// }

// // ClientSetUsername sets the client username
// func (action RoomAction) ClientSetUsername() {
// 	username := ""
// 	err := json.Unmarshal(action.Data, &username)
// 	if err != nil {
// 		log.Println("ClientSetUsername error:", err)
// 		return
// 	}
// 	action.Client.Username = username
// }

// // ClientSetColor sets the client color
// func (action RoomAction) ClientSetColor() {
// 	color := ""
// 	err := json.Unmarshal(action.Data, &color)
// 	if err != nil {
// 		log.Println("ClientSetColor error:", err)
// 		return
// 	}
// 	action.Client.Color = color
// }

// // ChatMessage broadcasts a new chat message
// func (action RoomAction) ChatMessage() {
// 	text := ""
// 	err := json.Unmarshal(action.Data, &text)
// 	if err != nil {
// 		log.Println("ChatMessage error:", err)
// 		return
// 	}
// 	action.Room.Chat.Messages = append(action.Room.Chat.Messages, ChatMessage{
// 		ID:        uuid.New(),
// 		Client:    action.Client,
// 		Timestamp: time.Now().UTC(),
// 		Text:      text,
// 	})
// }

// // PlayerInit inits the player
// func (action RoomAction) PlayerInit() {
// 	//
// }

// // PlayerTogglePlaying toggles playing state
// func (action RoomAction) PlayerTogglePlaying() {
// 	if action.Room.Player.CurrentVideo == nil {
// 		log.Println("PlayerTogglePlaying: CurrentVideo is nil")
// 		return
// 	}
// 	action.Room.Player.Playing = !action.Room.Player.Playing
// }

// // PlayerSkip skips to timestamp
// func (action RoomAction) PlayerSkip() {
// 	if len(action.Room.Player.Queue.Items) >= 1 {
// 		action.Room.Player.CurrentVideo = action.Room.Player.Queue.Items[0]
// 		action.Room.Player.Queue.RemoveItem(action.Room.Player.CurrentVideo.ID)
// 		log.Println("PlayerSkip: Skipped by:", action.Client.ID)
// 		return
// 	}
// 	log.Println("PlayerSkip: Queue is empty")
// }

// // QueueAdd adds a queue item to the queue
// func (action RoomAction) QueueAdd() {
// 	video := Video{}
// 	err := json.Unmarshal(action.Data, &video)
// 	if err != nil {
// 		log.Println("QueueAdd error:", err)
// 		return
// 	}
// 	if action.Room.Player.CurrentVideo == nil {
// 		action.Room.Player.Play(action.Client.NewQueueItem(video))
// 		return
// 	}
// 	action.Client.AddVideoToQueue(video)
// }

// // Play video
// func (p *Player) Play(queueItem *QueueItem) {
// 	p.CurrentVideo = queueItem
// 	p.Playing = true
// }

// // AddVideoToQueue adds a new video to the queue
// func (c *Client) AddVideoToQueue(video Video) {
// 	c.Room.Player.Queue.Items = append(c.Room.Player.Queue.Items, c.NewQueueItem(video))
// }

// // NewQueueItem returns a new queue item
// func (c *Client) NewQueueItem(video Video) *QueueItem {
// 	votes := make(map[uuid.UUID]*Vote)
// 	votes[c.ID] = &Vote{Client: c}
// 	return &QueueItem{
// 		ID:      uuid.New(),
// 		Video:   video,
// 		Votes:   votes,
// 		AddedBy: c,
// 	}
// }

// // QueueRemove removes a queue item from the queue
// func (action RoomAction) QueueRemove() {
// 	queueItemIDString := ""
// 	err := json.Unmarshal(action.Data, &queueItemIDString)
// 	if err != nil {
// 		log.Println("QueueRemove error:", err)
// 		return
// 	}
// 	queueItemUUID, _ := uuid.Parse(queueItemIDString)
// 	if queueItemUUID.String() != "" {
// 		action.Room.Player.Queue.RemoveItem(queueItemUUID)
// 	}
// }

// // RemoveItem from queue
// func (q *Queue) RemoveItem(id uuid.UUID) {
// 	for i, item := range q.Items {
// 		if item.ID == id {
// 			q.Items = append(q.Items[:i], q.Items[i+1:]...)
// 			return
// 		}
// 	}
// }

// // QueueVote toggles the vote on a queue item
// func (action RoomAction) QueueVote() {
// 	queueItemIDString := ""
// 	err := json.Unmarshal(action.Data, &queueItemIDString)
// 	if err != nil {
// 		log.Println("QueueVote error:", err)
// 		return
// 	}
// 	queueItemUUID, _ := uuid.Parse(queueItemIDString)
// 	if queueItemUUID.String() != "" {
// 		action.Client.Vote(queueItemUUID)
// 	}
// }
