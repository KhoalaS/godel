package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func NewDownloaderNode() Node {
	return Node{
		Type: "downloader",
		Run:  DownloaderNodeFunc,
		Io: map[string]*NodeIO{
			"job": {
				Type:      IOTypeGenerated,
				Id:        "job",
				ValueType: ValueTypeDownloadJob,
				Label:     "Downloader",
				Required:  true,
			},
			"url": {
				Type:      IOTypePassthrough,
				Id:        "url",
				ValueType: ValueTypeString,
				Label:     "Url",
				Required:  true,
			},
		},
		Name:     "Downloader",
		Status:   StatusPending,
		Category: NodeCategoryDownloader,
	}
}

func DownloaderNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	job := types.NewDownloadJob()

	jobUrl, ok := utils.FromAny[string](node.Io["url"].Value).Value()
	if !ok || jobUrl == "" {
		return NewInvalidNodeIOError(&node, "url")
	}

	job.Url = jobUrl
	node.Io["job"].Value = job

	return nil
}
