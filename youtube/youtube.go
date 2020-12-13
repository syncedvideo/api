package youtube

import (
	"net/http"
	"time"

	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

// YouTube is a youtube.Service wrapper
type YouTube struct {
	Service *youtube.Service
}

// New returns a new YouTube service
func New(apiKey string) *YouTube {
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
		Timeout:   time.Second * 5,
	}
	service, _ := youtube.New(client)
	return &YouTube{Service: service}
}

type videoSearchResults struct {
	query  string
	videos []*youTubeVideo
}

type youTubeVideo struct {
	ID             string
	Query          string
	Snippet        *youtube.SearchResultSnippet
	ContentDetails *youtube.VideoContentDetails
	Statistics     *youtube.VideoStatistics
}

// SearchVideos searches for videos
func (yt YouTube) SearchVideos(query string) ([]*youTubeVideo, error) {

	videos := []*youTubeVideo{}
	const maxResults = 10

	// get video ids and snippet
	searchListRequest := yt.Service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults)
	searchListResponse, err := searchListRequest.Do()
	if err != nil {
		return nil, err
	}

	for _, item := range searchListResponse.Items {
		if item.Id.VideoId != "" {
			videos = append(videos, &youTubeVideo{
				ID:      item.Id.VideoId,
				Snippet: item.Snippet,
			})
		}
	}

	// get video ids
	videoIDs := []string{}
	for _, video := range videos {
		videoIDs = append(videoIDs, video.ID)
	}

	// add content details to results
	videosListRequest := yt.Service.Videos.List([]string{"snippet", "contentDetails", "statistics"}).Id(videoIDs...)
	videosListResponse, err := videosListRequest.Do()
	if err != nil {
		return nil, err
	}

	for _, videoListItem := range videosListResponse.Items {
		for _, video := range videos {
			if videoListItem.Id == video.ID {
				video.ContentDetails = videoListItem.ContentDetails
				video.Statistics = videoListItem.Statistics
				// duration, _ := iso8601.ParseISO8601(videoListItem.ContentDetails.Duration)
				// result.Title = videoListItem.Snippet.Title
				// result.Description = videoListItem.Snippet.Description
				// result.Thumbnail = videoListItem.Snippet.Thumbnails.Default.Url
				// result.Duration = (duration.TM * 60) + duration.TS
				// result.ViewCount = videoListItem.Statistics.ViewCount
				// result.LikeCount = videoListItem.Statistics.LikeCount
				// result.DislikeCount = videoListItem.Statistics.DislikeCount
			}
		}
	}

	return videos, nil
}
