package syncedvideo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestPostRoom(t *testing.T) {

	store := &StubRoomStore{}
	server := NewServer(store, nil)

	t.Run("create room and return created room as json", func(t *testing.T) {
		wantRoom := Room{Name: "TestRoom"}

		request := NewPostRoomRequest(wantRoom)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertCreateRoomCalls(t, store.CreateRoomCalls, 1)

		gotRoom := GetRoomFromResponse(t, response.Body)
		AssertRoom(t, wantRoom, gotRoom)

		AssertStatus(t, response, http.StatusCreated)
		AssertJsonContentType(t, response)
	})
}

func TestGetRoom(t *testing.T) {

	jeromesRoom := Room{ID: "jerome", Name: "Jeromes room"}
	philippsRoom := Room{ID: "philipp", Name: "Philipps room"}
	store := &StubRoomStore{
		Rooms: map[string]Room{
			"jerome":  jeromesRoom,
			"philipp": philippsRoom,
		},
	}
	server := NewServer(store, nil)

	t.Run("return Jeromes room", func(t *testing.T) {
		request := NewGetRoomRequest("jerome")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetRoomFromResponse(t, response.Body)

		AssertRoom(t, jeromesRoom, got)
		AssertStatus(t, response, http.StatusOK)
		AssertJsonContentType(t, response)
	})

	t.Run("return Philipps room", func(t *testing.T) {
		request := NewGetRoomRequest("philipp")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetRoomFromResponse(t, response.Body)

		AssertRoom(t, philippsRoom, got)
		AssertStatus(t, response, http.StatusOK)
		AssertJsonContentType(t, response)
	})

	t.Run("return 404 on missing rooms", func(t *testing.T) {
		request := NewGetRoomRequest("tobi")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusNotFound)
	})
}

func TestWebSocket(t *testing.T) {

	store := &StubRoomStore{
		Rooms: map[string]Room{
			"jerome": {
				ID:   "jerome",
				Name: "Jeromes room",
			},
		},
	}

	pubSub := &MockPubSub{}
	server := NewServer(store, pubSub)

	t.Run("cannot establish websocket connection if room is missing", func(t *testing.T) {

		wsServer := httptest.NewServer(server)
		defer wsServer.Close()

		wsURL := newWebSocketURL(wsServer.URL, "philipp")
		_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

		AssertError(t, websocket.ErrBadHandshake, err)
	})
}

func TestPostChat(t *testing.T) {

	store := &StubRoomStore{
		Rooms: map[string]Room{
			"jerome": {
				ID:   "jerome",
				Name: "Jeromes room",
			},
		},
	}

	pubSub := &MockPubSub{}
	server := NewServer(store, pubSub)

	t.Run("return 404 on missing rooms", func(t *testing.T) {
		request := NewPostRoomChatRequest("philipp", ChatMessage{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusNotFound)
	})

	t.Run("send and receive chat message", func(t *testing.T) {
		wsServer := httptest.NewServer(server)
		wsURL := newWebSocketURL(wsServer.URL, "jerome")
		ws := MustDialWS(t, wsURL)

		defer ws.Close()
		defer wsServer.Close()

		chatMsg := NewChatMessage("Tobi", "Steinreinigung läuft")

		request := NewPostRoomChatRequest("jerome", chatMsg)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusCreated)
		AssertJsonContentType(t, response)

		chatMsgB, _ := json.Marshal(chatMsg)
		wantEvent := NewEvent(EventChat, chatMsgB)
		Within(t, 10*time.Millisecond, func() {
			AssertWebsocketGotEvent(t, ws, wantEvent)
		})
	})
}

func newWebSocketURL(serverURL, roomID string) string {
	return "ws" + strings.TrimPrefix(serverURL, "http") + fmt.Sprintf("/rooms/%s/ws", roomID)
}
