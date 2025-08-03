package utils

import (
	"net/http"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	err := Download(client, "http://localhost:8080/files/test.txt", "", nil)

	if err != nil {
		t.Error(err)
	}
}
