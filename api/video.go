package syncedvideo

import (
	"time"

	"github.com/google/uuid"
	iso8601 "github.com/senseyeio/duration"
	"github.com/syncedvideo/syncedvideo/youtube"
)

const VideoProviderYouTube = 1

type Video struct {
	ID         uuid.UUID       `json:"id"`
	Provider   int             `json:"provider"`
	ProviderID string          `json:"providerId"`
	Title      string          `json:"title"`
	Author     string          `json:"author"`
	Thumbnail  string          `json:"thumbnail"`
	Statistics videoStatistics `json:"statistics"`
	Duration   time.Duration   `json:"-"`
}

type videoStatistics struct {
	Views             uint64        `json:"views"`
	Likes             uint64        `json:"likes"`
	Dislikes          uint64        `json:"dislikes"`
	DurationInSeconds time.Duration `json:"durationInSeconds"`
}

func NewVideo(ytVideo youtube.Video) Video {
	d, _ := iso8601.ParseISO8601(ytVideo.ContentDetails.Duration)
	duration := d.TS
	duration += d.TM * 60
	duration += d.TH * (60 * 60)
	video := Video{
		ID:         uuid.New(),
		Provider:   VideoProviderYouTube,
		ProviderID: ytVideo.ID,
		Title:      ytVideo.Snippet.Title,
		Author:     ytVideo.Snippet.ChannelTitle,
		Thumbnail:  ytVideo.Snippet.Thumbnails.Default.Url,
	}
	video.Statistics = videoStatistics{
		Views:             ytVideo.Statistics.ViewCount,
		Likes:             ytVideo.Statistics.LikeCount,
		Dislikes:          ytVideo.Statistics.DislikeCount,
		DurationInSeconds: time.Duration(duration),
	}
	return video
}
