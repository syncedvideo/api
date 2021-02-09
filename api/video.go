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
	Views             uint64  `json:"views"`
	Likes             uint64  `json:"likes"`
	Dislikes          uint64  `json:"dislikes"`
	DurationInSeconds float64 `json:"durationInSeconds"`
}

func NewVideo(ytVideo youtube.Video) Video {
	d := parseISO8601(ytVideo.ContentDetails.Duration)
	video := Video{
		ID:         uuid.New(),
		Provider:   VideoProviderYouTube,
		ProviderID: ytVideo.ID,
		Title:      ytVideo.Snippet.Title,
		Author:     ytVideo.Snippet.ChannelTitle,
		Thumbnail:  ytVideo.Snippet.Thumbnails.Default.Url,
		Duration:   d,
	}
	video.Statistics = videoStatistics{
		Views:             ytVideo.Statistics.ViewCount,
		Likes:             ytVideo.Statistics.LikeCount,
		Dislikes:          ytVideo.Statistics.DislikeCount,
		DurationInSeconds: video.Duration.Seconds(),
	}
	return video
}

func parseISO8601(from string) time.Duration {
	d, _ := iso8601.ParseISO8601(from)
	duration := time.Duration(d.TS) * time.Second
	duration += time.Duration(d.TM) * time.Minute
	duration += time.Duration(d.TH) * time.Hour
	return duration
}
