package youtube

import (
	"net/http"
	"time"

	"github.com/google/uuid"
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
	videos []*YouTubeVideo
}

type YouTubeVideo struct {
	ID             uuid.UUID
	YouTubeID      string
	Query          string
	Snippet        *youtube.SearchResultSnippet
	ContentDetails *youtube.VideoContentDetails
	Statistics     *youtube.VideoStatistics
}

// SearchVideos searches for videos
func (yt YouTube) SearchVideos(query string) ([]*YouTubeVideo, error) {

	videos := []*YouTubeVideo{}
	const maxResults = 10

	// get video ids and snippet
	searchListRequest := yt.Service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults)
	searchListResponse, err := searchListRequest.Do()
	if err != nil {
		return nil, err
	}

	if len(searchListResponse.Items) == 0 {
		return videos, nil
	}

	for _, item := range searchListResponse.Items {
		if item.Id.VideoId != "" {
			videos = append(videos, &YouTubeVideo{
				ID:        uuid.New(),
				YouTubeID: item.Id.VideoId,
				Snippet:   item.Snippet,
			})
		}
	}

	// get video ids
	videoIDs := []string{}
	for _, video := range videos {
		videoIDs = append(videoIDs, video.YouTubeID)
	}

	// add content details to results
	videosListRequest := yt.Service.Videos.List([]string{"id", "snippet", "contentDetails", "statistics"}).Id(videoIDs...)
	videosListResponse, err := videosListRequest.Do()
	if err != nil {
		return nil, err
	}

	for _, videoListItem := range videosListResponse.Items {
		for _, video := range videos {
			if videoListItem.Id == video.YouTubeID {
				video.ContentDetails = videoListItem.ContentDetails
				video.Statistics = videoListItem.Statistics
			}
		}
	}

	return videos, nil
}
