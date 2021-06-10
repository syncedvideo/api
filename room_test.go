package syncedvideo

import (
	"testing"
)

func TestRoom(t *testing.T) {

	t.Run("play video", func(t *testing.T) {

		eventManager := NewMockEventManager()
		room := NewRoom(eventManager)
		video := &Video{
			ID:         "test",
			Provider:   "youtube",
			ProviderID: "yt-test",
			Title:      "Test video",
		}

		room.Play(video)

		gotEvent := eventManager.Events[0]

		AssertVideoPlayerIsPlaying(t, room)
		AssertEventType(t, EventPlay, gotEvent)
		AssertEventData(t, video, gotEvent)
	})

	t.Run("pause video", func(t *testing.T) {
		//
	})

	t.Run("seek video", func(t *testing.T) {
		//
	})

	t.Run("skip video", func(t *testing.T) {
		//
	})

	t.Run("send chat messsage", func(t *testing.T) {

		eventManager := NewMockEventManager()
		room := NewRoom(eventManager)

		chatMessage := NewChatMessage("Jerome", "Steinreinigung läuft")
		room.SendChatMessage(chatMessage)

		got := eventManager.Events[0]

		AssertEventType(t, EventChat, got)
		AssertEventData(t, chatMessage, got)
	})
}
