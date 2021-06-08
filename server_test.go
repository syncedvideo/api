package syncedvideo

import (
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
		AssertStatus(t, response.Code, http.StatusCreated)
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
		AssertStatus(t, response.Code, http.StatusOK)
		AssertJsonContentType(t, response)
		AssertCookie(t, response.Result().Cookies(), userCookieName)
	})

	t.Run("return Philipps room", func(t *testing.T) {
		request := NewGetRoomRequest("philipp")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetRoomFromResponse(t, response.Body)

		AssertRoom(t, philippsRoom, got)
		AssertStatus(t, response.Code, http.StatusOK)
		AssertJsonContentType(t, response)
		AssertCookie(t, response.Result().Cookies(), userCookieName)
	})

	t.Run("return 404 on missing rooms", func(t *testing.T) {
		request := NewGetRoomRequest("tobi")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusNotFound)
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

	eventManager := &MockEventManager{}
	server := NewServer(store, eventManager)

	requestHeader := make(http.Header)
	requestHeader.Add("Cookie", NewUserCookie().String())

	t.Run("fail if unauthorized", func(t *testing.T) {

		wsServer := httptest.NewServer(server)
		defer wsServer.Close()

		wsURL := newWebSocketURL(wsServer.URL, "jerome")
		_, response, err := websocket.DefaultDialer.Dial(wsURL, nil)

		AssertError(t, websocket.ErrBadHandshake, err)
		AssertStatus(t, response.StatusCode, http.StatusUnauthorized)
	})

	t.Run("fail if room is missing", func(t *testing.T) {

		wsServer := httptest.NewServer(server)
		defer wsServer.Close()

		wsURL := newWebSocketURL(wsServer.URL, "philipp")
		_, response, err := websocket.DefaultDialer.Dial(wsURL, requestHeader)

		AssertError(t, websocket.ErrBadHandshake, err)
		AssertStatus(t, response.StatusCode, http.StatusNotFound)
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

	eventManager := &MockEventManager{}
	server := NewServer(store, eventManager)

	t.Run("return 401 if unauthorized", func(t *testing.T) {
		request := NewPostRoomChatRequest("jerome", ChatMessage{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("return 404 on missing rooms", func(t *testing.T) {
		request := NewPostRoomChatRequest("philipp", ChatMessage{})
		request.AddCookie(NewUserCookie())

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("send and receive chat message", func(t *testing.T) {
		wsServer := httptest.NewServer(server)
		defer wsServer.Close()

		wsURL := newWebSocketURL(wsServer.URL, "jerome")
		requestHeader := make(http.Header)
		requestHeader.Add("Cookie", NewUserCookie().String())

		ws := MustDialWS(t, wsURL, requestHeader)
		defer ws.Close()

		chatMsg := NewChatMessage("Tobi", "Steinreinigung läuft")

		request := NewPostRoomChatRequest("jerome", chatMsg)
		request.AddCookie(NewUserCookie())

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusCreated)
		AssertJsonContentType(t, response)

		Within(t, 10*time.Millisecond, func() {
			AssertWebsocketGotEvent(t, ws, EventChat)
		})
	})
}

func newWebSocketURL(serverURL, roomID string) string {
	return "ws" + strings.TrimPrefix(serverURL, "http") + fmt.Sprintf("/rooms/%s/ws", roomID)
}
