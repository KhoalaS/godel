package nodes

import (
	"context"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
)

func CreateDownloaderNode() pipeline.Node {
	return pipeline.Node{
		Type: "downloader",
		Run:  DownloadNodeFunc,
		Io: map[string]*pipeline.NodeIO{
			"job": {
				Type:      pipeline.IOTypeOutput,
				Id:        "job",
				ValueType: pipeline.ValueTypeDownloadJob,
				Label:     "Downloader",
				Required:  true,
			},
			"url": {
				Type:      pipeline.IOTypePassthrough,
				Id:        "url",
				ValueType: pipeline.ValueTypeString,
				Label:     "Url",
				Required:  true,
			},
		},
		Name:     "Downloader",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.DownloaderNode,
	}
}

func DownloaderNodeFunc(ctx context.Context, node pipeline.Node, comm chan<- pipeline.PipelineMessage) error {
	job := types.NewDownloadJob()
	job.Url = (node.Io["url"].Value).(string)
	node.Io["job"].Value = job
	return nil
}
