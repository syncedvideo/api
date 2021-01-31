package middleware

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/request"
	"github.com/syncedvideo/syncedvideo/http/response"
)

func RoomMiddleware(next http.Handler, rs syncedvideo.RoomStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(chi.URLParam(r, "roomID"))
		if err != nil {
			log.Printf("error parsing uuid: %v", err)
			response.WithError(w, "room not found", http.StatusNotFound)
			return
		}
		room, err := rs.Get(id)
		if err == sql.ErrNoRows {
			response.WithError(w, "room not found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("error getting room: %v", err)
			response.WithError(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		request.WithRoom(r, room)
		next.ServeHTTP(w, r)
	})
}
