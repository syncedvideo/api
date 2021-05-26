package syncedvideo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
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

	t.Run("it returns Jeromes room", func(t *testing.T) {
		request := NewGetRoomRequest("jerome")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetRoomFromResponse(t, response.Body)

		AssertRoom(t, jeromesRoom, got)
		AssertStatus(t, response, http.StatusOK)
		AssertJsonContentType(t, response)
	})

	t.Run("it returns Philipps room", func(t *testing.T) {
		request := NewGetRoomRequest("philipp")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetRoomFromResponse(t, response.Body)

		AssertRoom(t, philippsRoom, got)
		AssertStatus(t, response, http.StatusOK)
		AssertJsonContentType(t, response)
	})

	t.Run("it returns 404 on missing rooms", func(t *testing.T) {
		request := NewGetRoomRequest("tobi")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusNotFound)
	})
}

func TestChat(t *testing.T) {

	roomID := "jerome"

	store := &StubRoomStore{
		Rooms: map[string]Room{
			roomID: {
				ID:   roomID,
				Name: "Jeromes room",
				Chat: &Chat{
					Messages: []ChatMessage{},
				}},
		},
	}

	pubSub := &MockRoomPubSub{}
	server := NewServer(store, pubSub)
	webSocketServer := httptest.NewServer(server)
	ws := mustDialWS(t, "ws"+strings.TrimPrefix(webSocketServer.URL, "http")+fmt.Sprintf("/rooms/%s/ws", roomID))

	defer ws.Close()
	defer webSocketServer.Close()

	t.Run("send and receive chat messages", func(t *testing.T) {
		wantChatMessage := ChatMessage{ID: "1", Author: "Tobi", Message: "Steinreinigung läuft"}
		request := NewPostRoomChatRequest(roomID, wantChatMessage)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusCreated)
		AssertJsonContentType(t, response)

		gotChatMessage := store.GetRoom(roomID).Chat.Messages[0]
		AssertChatMessage(t, gotChatMessage, wantChatMessage)

		messageBytes, _ := json.Marshal(&gotChatMessage)
		wantEvent := RoomEvent{T: 1, D: messageBytes}
		within(t, tenMs, func() {
			assertWebsocketGotEvent(t, ws, wantEvent)
		})
	})
}

var tenMs = 10 * time.Millisecond

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}

func mustDialWS(t testing.TB, wsURL string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	return ws
}

func mustWriteWSMessage(t testing.TB, conn *websocket.Conn, msg []byte) {
	t.Helper()
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func assertWebsocketGotEvent(t testing.TB, ws *websocket.Conn, want RoomEvent) {
	t.Helper()

	_, got, _ := ws.ReadMessage()

	wantBytes, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(got, wantBytes) {
		t.Errorf("got %q, want %q", got, wantBytes)
	}
}
