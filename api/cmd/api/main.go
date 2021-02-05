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
	apiHTTPPort         = os.Getenv("API_HTTP_PORT")
	apiPostgresHost     = os.Getenv("API_POSTGRES_HOST")
	apiPostgresPort     = os.Getenv("API_POSTGRES_PORT")
	apiPostgresDB       = os.Getenv("API_POSTGRES_DB")
	apiPostgresUser     = os.Getenv("API_POSTGRES_USER")
	apiPostgresPassword = os.Getenv("API_POSTGRES_PASSWORD")
	apiRedisHost        = os.Getenv("API_REDIS_HOST")
	apiRedisPort        = os.Getenv("API_REDIS_PORT")
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
	handler.RegisterUserHandler(router)
	handler.RegisterRoomHandler(router)

	// run http server
	log.Printf("http server listening on port %s\n", apiHTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", apiHTTPPort), router)
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
