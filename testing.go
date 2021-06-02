package syncedvideo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type StubRoomStore struct {
	Rooms           map[string]Room
	CreateRoomCalls []string
}

func (s *StubRoomStore) CreateRoom(room *Room) {
	s.CreateRoomCalls = append(s.CreateRoomCalls, room.Name)
}

func (s *StubRoomStore) GetRoom(id string) Room {
	return s.Rooms[id]
}

func NewPostRoomRequest(room Room) *http.Request {
	roomB, _ := json.Marshal(room)
	request, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(roomB))
	return request
}

func NewGetRoomRequest(id string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/rooms/%s", id), nil)
	return request
}

func NewPostRoomChatRequest(id string, message ChatMessage) *http.Request {
	messageB, _ := json.Marshal(message)
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/rooms/%s/chat", id), bytes.NewBuffer(messageB))
	return request
}

func AssertChatMessage(t testing.TB, got, want ChatMessage) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wrong message: got %v, want %v", got, want)
	}
}

func AssertStatus(t testing.TB, r *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := r.Code
	if got != want {
		t.Errorf("wrong status code: got %d, want %d", got, want)
	}
}

func AssertBody(t testing.TB, r *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := r.Body.String()
	if got != want {
		t.Errorf("wrong response body: got %q, want %q", got, want)
	}
}

func AssertCookie(t testing.TB, cookies []*http.Cookie, name string) {
	t.Helper()

	var got *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == name {
			got = cookie
			break
		}
	}
	if got == nil {
		t.Errorf(`cookie was not found in response, want "%s"`, name)
	}
}

func AssertNoCookie(t testing.TB, cookies []*http.Cookie, name string) {
	t.Helper()

	var got *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == name {
			got = cookie
			break
		}
	}
	if got != nil {
		t.Errorf("didn't expect an cookie but got one: %s", got.Name)
	}
}

func AssertJsonContentType(t testing.TB, r *httptest.ResponseRecorder) {
	t.Helper()
	got := r.Header().Get("Content-Type")
	want := jsonContentType
	if got != want {
		t.Errorf("wrong response content type: got %q, want %q", got, want)
	}
}

func AssertRoom(t testing.TB, want, got Room) {
	t.Helper()
	if want.Name != got.Name {
		t.Errorf("wrong room: got %s, want %s", got, want)

	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}
}

func AssertError(t testing.TB, want, got error) {
	t.Helper()
	if want != got {
		t.Fatalf("wrong error: got %s, want %s", got, want)
	}
}

func GetRoomFromResponse(t testing.TB, body io.Reader) Room {
	room, err := NewRoom(body)
	AssertNoError(t, err)
	return room
}

type MockPubSub struct {
	ch chan Event
}

func (m *MockPubSub) Publish(roomID string, event Event) {
	m.ch <- event
}

func (m *MockPubSub) Subscribe(roomID string) <-chan Event {
	ch := make(chan Event)
	m.ch = ch
	return ch
}

func Within(t testing.TB, d time.Duration, assert func()) {
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

func MustDialWS(t testing.TB, wsURL string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	return ws
}

func MustWriteWSMessage(t testing.TB, conn *websocket.Conn, msg []byte) {
	t.Helper()
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func AssertWebsocketGotEvent(t testing.TB, ws *websocket.Conn, want Event) {
	t.Helper()

	_, msg, err := ws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	got := Event{}
	err = json.Unmarshal(msg, &got)
	if err != nil {
		t.Fatal(err)
	}

	resetEventIDFields(t, &got)
	resetEventIDFields(t, &want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func resetEventIDFields(t testing.TB, event *Event) {
	t.Helper()

	event.ID = ""

	data := make(map[string]interface{})
	err := json.Unmarshal(event.D, &data)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := data["id"]
	if ok {
		data["id"] = ""
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	event.D = dataB
}

func AssertCreateRoomCalls(t testing.TB, got []string, want int) {
	t.Helper()
	if len(got) != want {
		t.Errorf("wrong create room calls: got %d, want %d", len(got), want)
	}
}
