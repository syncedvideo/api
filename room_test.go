package syncedvideo

import (
	"testing"
)

func TestRoom(t *testing.T) {

	t.Run("play video", func(t *testing.T) {

		eventManager := NewMockEventManager()
		room := NewRoom(eventManager)

		video := &Video{ID: "test", Name: "Test video"}
		room.PlayVideo(video)

		got := eventManager.Events[0]

		AssertEventType(t, EventPlayVideo, got)
		AssertEventData(t, video, got)
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
