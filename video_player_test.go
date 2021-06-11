package syncedvideo

import "testing"

// PLAY VIDEO
// set video
// start ticker -> set playtime on each tick

// PAUSE VIDEO
// stop ticker

// SEEK VIDEO
// stop ticker
// set playtime
// start ticker

// SKIP VIDEO
// set video
// start ticker

// TICKER
// countdown video duration
// set playtime on each tick
// play next video

func TestVideoPlayerTimer(t *testing.T) {
	spySleeper := &SpySleeper{}
	videoPlayer := &VideoPlayer{}
	videoPlayer.CurrentVideo = &Video{
		ID:         "test",
		Provider:   "youtube",
		ProviderID: "yt-test",
		Title:      "Test video",
		Duration:   5,
	}

	videoPlayer.Timer(spySleeper)

	AssertCurrentTime(t, videoPlayer.CurrentTime, videoPlayer.CurrentVideo.Duration)
	AssertSleeperCalls(t, spySleeper.Calls, videoPlayer.CurrentVideo.Duration+1)
}

func AssertCurrentTime(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("wrong current time: got %d, want %d", got, want)
	}
}

func AssertSleeperCalls(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("not enough calls to sleeper: got %d, want %d", got, want)
	}
}

type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}
