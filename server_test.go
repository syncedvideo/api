package syncedvideo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

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

func TestChat(t *testing.T) {

	store := &StubRoomStore{
		Rooms: map[string]Room{
			"jerome": {
				ID:   "jerome",
				Name: "Jeromes room",
			},
		},
	}

	pubSub := &MockRoomPubSub{}
	server := NewServer(store, pubSub)

	t.Run("return 404 on missing rooms", func(t *testing.T) {
		request := NewPostRoomChatRequest("philipp", ChatMessage{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusNotFound)
	})

	t.Run("send and receive chat message", func(t *testing.T) {
		wsServer := httptest.NewServer(server)
		ws := MustDialWS(t, "ws"+strings.TrimPrefix(wsServer.URL, "http")+fmt.Sprintf("/rooms/%s/ws", "jerome"))

		defer ws.Close()
		defer wsServer.Close()

		chatMsg := NewChatMessage("Tobi", "Steinreinigung läuft")

		request := NewPostRoomChatRequest("jerome", chatMsg)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusCreated)
		AssertJsonContentType(t, response)

		chatMsgB, _ := json.Marshal(chatMsg)
		wantEvent := NewRoomEvent(1, chatMsgB)
		Within(t, 10*time.Millisecond, func() {
			AssertWebsocketGotEvent(t, ws, wantEvent)
		})
	})
}
