package transformer

import (
	"os"
	"testing"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestRedditTransformer(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel)
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}

	link := os.Getenv("REDDIT_POST_LINK")

	job := types.DownloadJob{
		Url: link,
		Id:  "1",
	}

	err = RedditTransformer(&job)
	if err != nil {
		t.Error(err)
	}
}
