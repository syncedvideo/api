package syncedvideo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func AssertJSONContentType(t testing.TB, r *httptest.ResponseRecorder) {
	t.Helper()
	got := r.Header().Get("Content-Type")
	want := JSONContentType
	if got != want {
		t.Errorf("wrong response content type: got %q, want %q", got, want)
	}
}
