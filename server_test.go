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

		gotRoom := GetRoomFromResponse(t, response.Body)
		AssertRoom(t, wantRoom, gotRoom)
		AssertCreateRoomCalls(t, store.CreateRoomCalls, 1)
		AssertStatus(t, response, http.StatusCreated)
		AssertJsonContentType(t, response)
		AssertCookie(t, response.Result().Cookies(), userCookieName)
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
		AssertCookie(t, response.Result().Cookies(), userCookieName)
	})

	t.Run("return Philipps room", func(t *testing.T) {
		request := NewGetRoomRequest("philipp")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetRoomFromResponse(t, response.Body)

		AssertRoom(t, philippsRoom, got)
		AssertStatus(t, response, http.StatusOK)
		AssertJsonContentType(t, response)
		AssertCookie(t, response.Result().Cookies(), userCookieName)
	})

	t.Run("return 404 on missing rooms", func(t *testing.T) {
		request := NewGetRoomRequest("tobi")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusNotFound)
		AssertNoCookie(t, response.Result().Cookies(), userCookieName)
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

	// TODO: check user cookie
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

	t.Run("return 401 if unauthorized", func(t *testing.T) {
		request := NewPostRoomChatRequest("jerome", ChatMessage{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusUnauthorized)
	})

	t.Run("return 404 on missing rooms", func(t *testing.T) {
		request := NewPostRoomChatRequest("philipp", ChatMessage{})
		request.AddCookie(NewUserCookie())

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusNotFound)
	})

	t.Run("send and receive chat message", func(t *testing.T) {
		wsServer := httptest.NewServer(server)
		wsURL := newWebSocketURL(wsServer.URL, "jerome")
		ws := MustDialWS(t, wsURL)

		defer wsServer.Close()
		defer ws.Close()

		chatMsg := NewChatMessage("Tobi", "Steinreinigung läuft")

		request := NewPostRoomChatRequest("jerome", chatMsg)
		request.AddCookie(NewUserCookie())

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
