package pipeline

import (
	"context"
	"testing"

	"github.com/KhoalaS/godel/pkg/types"
)

func TestDownloaderNodeFunc(t *testing.T) {
	ctx := context.Background()

	node := CreateDownloaderNode()
	node.Io["url"].Value = "123"

	err := DownloaderNodeFunc(ctx, node, make(chan<- PipelineMessage), "1", node.Id)
	if err != nil {
		t.Fail()
	}

	if node.Io["job"].Value == nil {
		t.Fail()
	}

	switch node.Io["job"].Value.(type) {
	case *types.DownloadJob:
		t.Log("job value has correct type")
	default:
		t.Fail()
	}
}
