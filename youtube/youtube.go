package youtube

import (
	"net/http"
	"time"

	iso8601 "github.com/senseyeio/duration"
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

type VideoSearchResults struct {
	Query string               `json:"query"`
	Items []*VideoSearchResult `json:"items"`
}

func (results *VideoSearchResults) getItemIds() []string {
	ids := []string{}
	for _, result := range results.Items {
		ids = append(ids, result.ID)
	}
	return ids
}

type VideoSearchResult struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Duration     int    `json:"duration"`
	Thumbnail    string `json:"thumbnail"`
	ViewCount    uint64 `json:"viewCount"`
	LikeCount    uint64 `json:"likeCount"`
	DislikeCount uint64 `json:"dislikeCount"`
}

// VideoSearch searches for videos
func (yt YouTube) VideoSearch(query string) (*VideoSearchResults, error) {

	const maxResults = 10

	// get video ids and snippet
	searchListRequest := yt.Service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults)
	searchListResponse, err := searchListRequest.Do()
	if err != nil {
		return nil, err
	}

	results := &VideoSearchResults{
		Query: query,
		Items: []*VideoSearchResult{},
	}
	for _, item := range searchListResponse.Items {
		if item.Id.VideoId != "" {
			results.Items = append(results.Items, &VideoSearchResult{
				ID: item.Id.VideoId,
			})
		}
	}

	// add content details to results
	videosListRequest := yt.Service.Videos.List([]string{"snippet", "contentDetails", "statistics"}).Id(results.getItemIds()...)
	videosListResponse, err := videosListRequest.Do()
	if err != nil {
		return nil, err
	}

	for _, videoListItem := range videosListResponse.Items {
		for _, result := range results.Items {
			if videoListItem.Id == result.ID {
				duration, _ := iso8601.ParseISO8601(videoListItem.ContentDetails.Duration)
				result.Title = videoListItem.Snippet.Title
				result.Description = videoListItem.Snippet.Description
				result.Thumbnail = videoListItem.Snippet.Thumbnails.Default.Url
				result.Duration = (duration.TM * 60) + duration.TS
				result.ViewCount = videoListItem.Statistics.ViewCount
				result.LikeCount = videoListItem.Statistics.LikeCount
				result.DislikeCount = videoListItem.Statistics.DislikeCount
			}
		}
	}

	return results, nil
}
