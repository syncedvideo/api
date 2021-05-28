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
	Rooms map[string]Room
}

func (s *StubRoomStore) GetRoom(id string) Room {
	return s.Rooms[id]
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
	if !reflect.DeepEqual(want, got) {
		t.Errorf("wrong room: got %v, want %v", got, want)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
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

	resetEventIDFields(&got)
	resetEventIDFields(&want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func resetEventIDFields(event *Event) {
	event.ID = ""

	data := make(map[string]interface{})
	json.Unmarshal(event.D, &data)

	_, ok := data["id"]
	if ok {
		data["id"] = ""
	}

	dataB, _ := json.Marshal(data)
	event.D = dataB
}
