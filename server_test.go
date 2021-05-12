package syncedvideo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRoom(t *testing.T) {

	store := &StubRoomStore{
		Rooms: map[string]Room{
			"jerome":  {id: "jerome", name: "Jeromes room"},
			"philipp": {id: "philipp", name: "Philipps room"},
		},
	}
	server := NewServer(store)

	t.Run("it returns Jeromes room", func(t *testing.T) {
		request := NewGetRoomRequest("jerome")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
		AssertJSONContentType(t, response)
		AssertBody(t, response, "Jeromes room")
	})

	t.Run("it returns Philipps room", func(t *testing.T) {
		request := NewGetRoomRequest("philipp")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusOK)
		AssertJSONContentType(t, response)
		AssertBody(t, response, "Philipps room")
	})

	t.Run("it returns 404 on missing rooms", func(t *testing.T) {
		request := NewGetRoomRequest("tobi")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response, http.StatusNotFound)
	})
}
