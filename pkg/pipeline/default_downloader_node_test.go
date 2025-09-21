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

	var p IPipeline = &Pipeline{
		Id:   "1",
		Comm: make(chan PipelineMessage, 12),
	}

	err := DownloaderNodeFunc(ctx, node, p)
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
