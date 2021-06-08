package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/syncedvideo/syncedvideo"
)

func main() {
	store := &syncedvideo.StubRoomStore{
		Rooms: map[string]syncedvideo.Room{
			"jerome": {ID: "jerome", Name: "Jeromes room"},
		},
	}
	eventManager := syncedvideo.NewRedisEventManager(newRedisClient())

	server := syncedvideo.NewServer(store, eventManager)
	log.Fatal(http.ListenAndServe(":3000", server))
}

func newRedisClient() *redis.Client {
	redisOpts, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s", "redis", "6379"))
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisOpts)
	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	return redisClient
}
