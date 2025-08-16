package transformer

import (
	"os"
	"testing"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/joho/godotenv"
)

func TestGofileTransformer(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}

	link := os.Getenv("GOFILE_TEST_LINK")

	job := types.DownloadJob{
		Url: link,
	}

	err = GofileTransformer(&job)
	if err != nil {
		t.Error(err)
	}
}
