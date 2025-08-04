package transformer

import (
	"os"
	"testing"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/joho/godotenv"
)

func TestRealDebridTransformer(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}

	link := os.Getenv("RD_TEST_LINK")
	filename := os.Getenv("RD_TEST_FILENAME")

	job := types.DownloadJob{
		Url:      link,
		Filename: "",
		Id:       "100",
		Password: "",
	}
	newJob, err := RealDebridTransformer(job)
	if err != nil {
		t.Error(err)
	}

	if newJob.Filename != filename {
		t.Fail()
	}
}
