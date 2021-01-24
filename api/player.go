package syncedvideo

import (
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	Video   *Video      `json:"video"`
	Time    int64       `json:"time"`
	Playing bool        `json:"playing"`
	Queue   *VideoQueue `json:"queue"`
}

func NewVideoPlayer() *Player {
	return &Player{
		Video:   nil,
		Time:    0,
		Playing: false,
		Queue:   NewVideoQueue(),
	}
}

// Play sets the current video and playing state
func (player *Player) Play(video *Video) {
	if video != player.Video {
		player.Video = video
		player.Time = 0
	}
	player.Playing = true

	go func() {
		for {
			time.Sleep(1 * time.Second)
			if player.Time >= player.Video.Duration {
				log.Println("STOP!")
				player.Playing = false
				player.Video = nil
				player.Time = 0
				return
			}
			if !player.Playing {
				log.Println("STOP!")
				return
			}
			player.Time = player.Time + 1
			log.Println(player.Time)
		}
	}()

	log.Println("PLAY")
}

type VideoQueue struct {
	Videos []*Video `json:"videos"`
}

// Sort queue items by vote count
func (queue *VideoQueue) Sort() {
	sort.SliceStable(queue.Videos, func(i, j int) bool {
		return len(queue.Videos[i].Votes) > len(queue.Videos[j].Votes)
	})
}

func (queue *VideoQueue) Find(id uuid.UUID) *Video {
	for _, video := range queue.Videos {
		if video.ID == id {
			return video
		}
	}
	return nil
}

func (queue *VideoQueue) IsQueued(id uuid.UUID) bool {
	return queue.Find(id) != nil
}

func (queue *VideoQueue) Add(user *User, video *Video) {
	if !queue.IsQueued(video.ID) {
		video.AddedBy = user
		queue.Videos = append(queue.Videos, video)
	} else {
		log.Println("Video is already queued:", video)
	}
}

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

func (queue *VideoQueue) ToggleVote(user *User, video *Video) {
	if !queue.IsQueued(video.ID) {
		return
	}
	video.ToggleVote(user)
}

func NewVideoQueue() *VideoQueue {
	return &VideoQueue{
		Videos: []*Video{},
	}
}
