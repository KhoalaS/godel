package pipeline

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/KhoalaS/godel/pkg/types"
)

func TestDownload(t *testing.T) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}

	ctx := context.TODO()

	job := types.DownloadJob{
		Url: "http://localhost:8080/files/test.txt",
	}

	var p IPipeline = &Pipeline{
		Id:   "1",
		Comm: make(chan PipelineMessage, 12),
	}

	err := Download(ctx, client, &job, p, "1")

	if err != nil {
		t.Error(err)
	}
}
