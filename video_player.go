package syncedvideo

type VideoPlayer struct {
	CurrentVideo *Video
	CurrentTime  int
}

type Video struct {
	ID         string `json:"id"`
	Provider   string `json:"provider"`
	ProviderID string `json:"providerId"`
	Title      string `json:"title"`
	Duration   int    `json:"duration"`
}

func (p *VideoPlayer) Timer(sleeper Sleeper) {
	for i := p.CurrentTime; i <= p.CurrentVideo.Duration; i++ {
		sleeper.Sleep()
		p.CurrentTime = i
	}
}
