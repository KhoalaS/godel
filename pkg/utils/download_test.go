package utils

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

	err := Download(ctx, client, &job, nil)

	if err != nil {
		t.Error(err)
	}
}
