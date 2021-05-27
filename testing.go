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

type MockRoomPubSub struct {
	ch chan RoomEvent
}

func (m *MockRoomPubSub) Publish(roomID string, event RoomEvent) {
	m.ch <- event
}

func (m *MockRoomPubSub) Subscribe(roomID string) <-chan RoomEvent {
	ch := make(chan RoomEvent)
	m.ch = ch
	return ch
}
