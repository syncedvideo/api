package syncedvideo

import (
	"github.com/google/uuid"
)

// Provider to identify the video provider
type Provider string

// YouTubeProvider for Youtube videos
const YouTubeProvider = Provider("youtube")

// Video represents a video that can be
// played by a Player or added to a VideoQueue
type Video struct {
	ID          uuid.UUID           `json:"id"`
	ProviderID  string              `json:"providerId"`
	Provider    Provider            `json:"provider"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Duration    int64               `json:"duration"`
	Thumbnail   string              `json:"thumbnail"`
	AddedBy     *User               `json:"addedBy"`
	Votes       map[uuid.UUID]*User `json:"votes"`
	Statistics  videoStatistics     `json:"statistics"`
}

type videoStatistics struct {
	ViewCount    uint64 `json:"viewCount"`
	LikeCount    uint64 `json:"likeCount"`
	DislikeCount uint64 `json:"dislikeCount"`
}

func (v *Video) AddVote(user *User) {
	v.Votes[user.ID] = user
}

func (v *Video) RemoveVote(user *User) {
	delete(v.Votes, user.ID)
}

func (v *Video) ToggleVote(user *User) {
	_, voted := v.Votes[user.ID]
	if voted {
		v.RemoveVote(user)
		return
	}
	v.AddVote(user)
}
