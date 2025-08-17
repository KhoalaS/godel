package transformer

import (
	"os"
	"testing"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/joho/godotenv"
)

func TestJpgfishTransformer(t *testing.T) {

	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}

	link := os.Getenv("JPGFISH_ALBUM_TEST_LINK")

	job := types.DownloadJob{
		Url: link,
	}

	err = JpgfishTransformer(&job)

	if err != nil {
		t.Fail()
	}
}
