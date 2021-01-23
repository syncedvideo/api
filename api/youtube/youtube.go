package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

type VideoInfo struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	LengthInSeconds int    `json:"lengthInSeconds"`
	Thumbnail       string `json:"thumbnail"`
	ViewCount       int    `json:"viewCount"`
	Author          string `json:"author"`
}

type playerResponseData struct {
	VideoDetails struct {
		VideoID         string `json:"videoId"`
		Title           string `json:"title"`
		LengthInSeconds string `json:"lengthSeconds"`
		Thumbnail       struct {
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		ViewCount string `json:"viewCount"`
		Author    string `json:"author"`
	} `json:"videoDetails"`
}

// GetVideoInfo downloads YT get_video_info endpoint
func GetVideoInfo(videoID string) (*VideoInfo, error) {
	u, err := url.Parse(fmt.Sprintf("https://youtube.com/get_video_info?video_id=%s", videoID))
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(bodyBytes)
	values, err := url.ParseQuery(bodyString)
	if err != nil {
		return nil, err
	}

	playerResponse := values.Get("player_response")
	var prData playerResponseData
	err = json.Unmarshal([]byte(playerResponse), &prData)
	if err != nil {
		return nil, err
	}

	lengthInSeconds, err := strconv.Atoi(prData.VideoDetails.LengthInSeconds)
	if err != nil {
		return nil, err
	}

	viewCount, err := strconv.Atoi(prData.VideoDetails.ViewCount)
	if err != nil {
		return nil, err
	}

	thumbnail := ""
	for _, t := range prData.VideoDetails.Thumbnail.Thumbnails {
		if t.Height >= 100 {
			thumbnail = t.URL
			break
		}
	}

	return &VideoInfo{
		ID:              prData.VideoDetails.VideoID,
		Title:           prData.VideoDetails.Title,
		LengthInSeconds: lengthInSeconds,
		Thumbnail:       thumbnail,
		ViewCount:       viewCount,
		Author:          prData.VideoDetails.Author,
	}, nil
}

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
