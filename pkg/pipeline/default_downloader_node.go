package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/types"
)

func CreateDownloaderNode() Node {
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

func DownloaderNodeFunc(ctx context.Context, node Node, comm chan<- PipelineMessage, pipelineId string, nodeId string) error {
	job := types.NewDownloadJob()

	job.Url = (node.Io["url"].Value).(string)
	node.Io["job"].Value = job

	return nil
}
