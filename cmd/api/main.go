package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/syncedvideo/syncedvideo"
	"github.com/syncedvideo/syncedvideo/http/handler"
	"github.com/syncedvideo/syncedvideo/http/middleware"
	"github.com/syncedvideo/syncedvideo/store/postgres"
)

var (
	apiHTTPPort         = os.Getenv("HTTP_PORT")
	apiPostgresHost     = os.Getenv("POSTGRES_HOST")
	apiPostgresPort     = os.Getenv("POSTGRES_PORT")
	apiPostgresDB       = os.Getenv("POSTGRES_DB")
	apiPostgresUser     = os.Getenv("POSTGRES_USER")
	apiPostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	apiRedisHost        = os.Getenv("REDIS_HOST")
	apiRedisPort        = os.Getenv("REDIS_PORT")
	ytAPIKey            = os.Getenv("YOUTUBE_API_KEY")
)

func main() {
	flag.Parse()

	// init store
	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", apiPostgresHost, apiPostgresUser, apiPostgresPassword, apiPostgresDB)
	store, err := postgres.NewStore(postgresDsn)
	if err != nil {
		panic(err)
	}

	// init redis client
	redisOpts, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s", apiRedisHost, apiRedisPort))
	if err != nil {
		panic(err)
	}
	redis := redis.NewClient(redisOpts)
	_, err = redis.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	// register config
	syncedvideo.RegisterConfig(store, redis)

	// register http handlers
	router := chi.NewRouter()
	router.Use(middleware.CorsMiddleware)
	handler.RegisterAuthHandler(router)
	handler.RegisterUserHandler(router)
	handler.RegisterRoomHandler(router, ytAPIKey)

	// run http server
	log.Printf("http server listening on port %s\n", apiHTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", apiHTTPPort), router)
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
