package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
