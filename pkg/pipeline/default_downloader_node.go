package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

func CreateDownloaderNode() Node {
	return Node{
		Type: "downloader",
		Run:  DownloaderNodeFunc,
		Io: map[string]*NodeIO{
			"job": {
				Type:      IOTypeOutput,
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
		NodeType: DownloaderNode,
	}
}

func DownloaderNodeFunc(ctx context.Context, node Node, comm chan<- PipelineMessage) error {
	job := types.NewDownloadJob()
	log.Debug().Any("url", (node.Io["url"].Value).(string)).Send()

	job.Url = (node.Io["url"].Value).(string)
	node.Io["job"].Value = job

	log.Debug().Any("downloader", node.Io["job"].Value).Send()
	return nil
}
