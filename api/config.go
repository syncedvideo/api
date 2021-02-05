package syncedvideo

import (
	"github.com/go-redis/redis/v8"
)

var Config config

type config struct {
	Store Store
	Redis *redis.Client
}

func RegisterConfig(store Store, redis *redis.Client) {
	Config = config{
		store,
		redis,
	}
}
