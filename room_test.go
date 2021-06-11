package syncedvideo

import (
	"testing"
)

func TestRoom(t *testing.T) {

	t.Run("send chat messsage", func(t *testing.T) {

		eventManager := NewMockEventManager()
		room := NewRoom(eventManager)

		chatMessage := NewChatMessage("Jerome", "Steinreinigung läuft")
		room.SendChatMessage(chatMessage)

		event := eventManager.Events[0]

		AssertEventType(t, event, EventChat)
		AssertEventData(t, event, chatMessage)
	})
}

// func TestRoomVideoPlayer(t *testing.T) {

// 	t.Run("play video", func(t *testing.T) {

// 		eventManager := NewMockEventManager()
// 		room := NewRoom(eventManager)
// 		room.Play(&Video{
// 			ID:         "test",
// 			Provider:   "youtube",
// 			ProviderID: "yt-test",
// 			Title:      "Test video",
// 		})

// 		event := eventManager.Events[0]

// 		AssertVideoPlayerIsPlaying(t, room)
// 		AssertEventType(t, event, EventPlay)
// 		AssertEventData(t, event, room.VideoPlayer)
// 	})

// 	t.Run("pause video", func(t *testing.T) {

// 		eventManager := NewMockEventManager()
// 		room := NewRoom(eventManager)

// 		room.Pause()

// 		event := eventManager.Events[0]

// 		AssertVideoPlayerIsPaused(t, room)
// 		AssertEventType(t, event, EventPause)
// 		AssertEventData(t, event, room.VideoPlayer)
// 	})

// 	t.Run("set playtime", func(t *testing.T) {

// 		eventManager := NewMockEventManager()
// 		room := NewRoom(eventManager)
// 		seconds := 10

// 		room.Playtime(seconds)

// 		event := eventManager.Events[0]

// 		if room.VideoPlayer.CurrentPlaytime() != seconds {
// 			t.Errorf("wrong current time: got %d, want %d", room.VideoPlayer.CurrentPlaytime(), seconds)
// 		}

// 		// AssertVideoPlayerCurrentTime(t, seconds, room)
// 		AssertEventType(t, event, EventSeek)
// 		AssertEventData(t, event, room.VideoPlayer)
// 	})

// 	t.Run("skip to next video", func(t *testing.T) {
// 		//
// 	})
// }
