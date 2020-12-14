package room

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// VideoPlayer represents the room's video player
type VideoPlayer struct {
	CurrentVideo *Video        `json:"currentVideo"`
	CurrentTime  time.Duration `json:"currentTime"`
	Playing      bool          `json:"playing"`
	Queue        *VideoQueue   `json:"queue"`
}

// NewVideoPlayer returns a new video player
func NewVideoPlayer() *VideoPlayer {
	return &VideoPlayer{
		CurrentVideo: nil,
		Playing:      false,
		Queue:        NewVideoQueue(),
	}
}

// Play sets the current video and playing state
func (player *VideoPlayer) Play(video *Video) {
	player.CurrentVideo = video
	player.Playing = true
}

// VideoQueue represents the room's video queue
type VideoQueue struct {
	Videos []*Video `json:"videos"`
}

// Find video in queue
func (queue *VideoQueue) Find(id uuid.UUID) *Video {
	for _, video := range queue.Videos {
		if video.ID == id {
			return video
		}
	}
	return nil
}

// IsQueued checks if video is queued
func (queue *VideoQueue) IsQueued(id uuid.UUID) bool {
	return queue.Find(id) != nil
}

// Add video to queue
func (queue *VideoQueue) Add(user *User, video *Video) {
	if !queue.IsQueued(video.ID) {
		video.AddedBy = user
		queue.Videos = append(queue.Videos, video)
	} else {
		log.Println("Video is already queued:", video)
	}
}

// Remove video from queue
func (queue *VideoQueue) Remove(id uuid.UUID) {
	if queue.IsQueued(id) {
		for i, video := range queue.Videos {
			if video.ID == id {
				queue.Videos = append(queue.Videos[:i], queue.Videos[i+1:]...)
				log.Println("Removed video from queue:", video)
				break
			}
		}
	}
}

// ToggleVote of queued video
func (queue *VideoQueue) ToggleVote(user *User, video *Video) {
	if !queue.IsQueued(video.ID) {
		return
	}
	video.ToggleVote(user)
}

// NewVideoQueue returns a new video queue
func NewVideoQueue() *VideoQueue {
	return &VideoQueue{
		Videos: []*Video{},
	}
}
