package room

import (
	"github.com/google/uuid"
)

// Provider to identify the video provider
type Provider string

// YouTubeProvider for Youtube videos
const YouTubeProvider = Provider("youtube")

// Video represents a video that can be
// played by a VideoPlayer or added to a VideoQueue
type Video struct {
	ID          uuid.UUID           `json:"id"`
	Provider    Provider            `json:"provider"`
	ProviderID  string              `json:"providerId"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Duration    int                 `json:"duration"`
	Thumbnail   string              `json:"thumbnail"`
	AddedBy     *User               `json:"addedBy"`
	Votes       map[uuid.UUID]*User `json:"votes"`
	Statistics  struct {
		ViewCount    uint64 `json:"viewCount"`
		LikeCount    uint64 `json:"likeCount"`
		DislikeCount uint64 `json:"dislikeCount"`
	} `json:"statistics"`
}

// ToggleVote of video
func (v *Video) ToggleVote(user *User) {
	_, voted := v.Votes[user.ID]
	if voted {
		delete(v.Votes, user.ID)
		return
	}
	v.Votes[user.ID] = user
}
