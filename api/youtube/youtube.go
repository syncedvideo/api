package youtube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/api/googleapi/transport"
	yt "google.golang.org/api/youtube/v3"
)

var ErrNoResults = errors.New("No results")

type YouTube struct {
	Service *yt.Service
}

func New(apiKey string) *YouTube {
	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
		Timeout:   time.Second * 5,
	}
	service, _ := yt.New(client)
	return &YouTube{Service: service}
}

type Video struct {
	ID             string                  `json:"id"`
	Snippet        *yt.VideoSnippet        `json:"snippet"`
	ContentDetails *yt.VideoContentDetails `json:"contentDetails"`
	Statistics     *yt.VideoStatistics     `json:"statistics"`
}

// MarshalBinary: Implementation of encoding.BinaryMarshaler interface
func (v Video) MarshalBinary() (data []byte, err error) {
	return json.Marshal(v)
}

func (yt YouTube) GetVideo(videoURL string) (Video, error) {
	videoID := ExtractVideoID(videoURL)
	if videoID == "" {
		return Video{}, errors.New("videoID is empty")
	}

	req := yt.Service.Videos.List([]string{"id", "snippet", "contentDetails", "statistics"}).Id(videoID)
	resp, err := req.Do()
	if err != nil {
		return Video{}, fmt.Errorf("Do failed: %w\n", err)
	}
	if len(resp.Items) == 0 {
		return Video{}, ErrNoResults
	}

	return Video{
		ID:             resp.Items[0].Id,
		Snippet:        resp.Items[0].Snippet,
		ContentDetails: resp.Items[0].ContentDetails,
		Statistics:     resp.Items[0].Statistics,
	}, nil
}

func ExtractVideoID(videoURL string) string {
	url, _ := url.Parse(videoURL)
	if url.Host == "" {
		return videoURL
	}
	return url.Query().Get("v")
}

func cacheKey(videoID string) string {
	return "video." + videoID
}

func GetVideoFromCache(r *redis.Client, videoID string) (Video, error) {
	res, err := r.Get(context.Background(), cacheKey(videoID)).Result()
	if err == redis.Nil {
		return Video{}, redis.Nil
	}
	if err != nil {
		return Video{}, fmt.Errorf("Redis Get failed: %w\n", err)
	}
	video := Video{}
	err = json.Unmarshal([]byte(res), &video)
	if err != nil {
		return Video{}, fmt.Errorf("Unmarshal failed: %w\n", err)
	}
	return video, nil
}

func CacheVideo(r *redis.Client, video Video) error {
	return r.Set(context.Background(), cacheKey(video.ID), video, 6*time.Hour).Err()
}
