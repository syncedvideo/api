package config

import "github.com/go-redis/redis/v8"

// Redis client
var Redis *redis.Client

// NewRedisClient inits the redis client
func NewRedisClient(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	Redis = redis.NewClient(opt)
	return Redis
}
