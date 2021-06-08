package syncedvideo

import (
	"testing"
)

func TestVideoPlayer(t *testing.T) {

	t.Run("add videos", func(t *testing.T) {

		room := Room{ID: "test", Name: "Test room"}
		video1 := &Video{ID: "video1", Name: "Test video 1"}
		video2 := &Video{ID: "video2", Name: "Test video 2"}

		room.AddVideo(video1)
		room.AddVideo(video2)

		AssertVideo(t, video1, room.Video)
		AssertVideo(t, video2, room.Playlist[0])
	})
}
