package postgres

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/syncedvideo/syncedvideo"
)

var (
	apiPostgresHost     = os.Getenv("POSTGRES_HOST")
	apiPostgresPort     = os.Getenv("POSTGRES_PORT")
	apiPostgresDB       = os.Getenv("POSTGRES_DB")
	apiPostgresUser     = os.Getenv("POSTGRES_USER")
	apiPostgresPassword = os.Getenv("POSTGRES_PASSWORD")
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
