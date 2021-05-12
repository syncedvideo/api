package syncedvideo

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
	server := NewServer(store)

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
