package main

import (
	"log"
	"net/http"

	"github.com/syncedvideo/syncedvideo"
)

func main() {
	store := syncedvideo.StubRoomStore{
		Rooms: map[string]syncedvideo.Room{
			"jerome": {ID: "jerome", Name: "Jeromes room"},
		},
	}
	server := syncedvideo.NewServer(&store)
	log.Fatal(http.ListenAndServe(":3000", server))
}
