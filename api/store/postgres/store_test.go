package postgres

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/syncedvideo/syncedvideo"
)

var (
	apiPostgresHost     = os.Getenv("API_POSTGRES_HOST")
	apiPostgresPort     = os.Getenv("API_POSTGRES_PORT")
	apiPostgresDB       = os.Getenv("API_POSTGRES_DB")
	apiPostgresUser     = os.Getenv("API_POSTGRES_USER")
	apiPostgresPassword = os.Getenv("API_POSTGRES_PASSWORD")
)

var store syncedvideo.Store

func TestMain(m *testing.M) {
	flag.Parse()
	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", apiPostgresHost, apiPostgresUser, apiPostgresPassword, apiPostgresDB)
	s, err := NewStore(postgresDsn)
	if err != nil {
		panic(err)
	}
	store = s
	os.Exit(m.Run())
}
