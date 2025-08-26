package nodes

import (
	"context"
	"net/http"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func CreateDownloaderNode() pipeline.Node {
	inputs := []pipeline.NodeInput{}
	return pipeline.Node{
		Type:     "default_downloader",
		Phase:    pipeline.DownloadPhase,
		Run:      DownloaderNodeFunc,
		Inputs:   inputs,
		Name:     "Standard Downloader",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.DownloadNode,
	}
}

func DownloaderNodeFunc(ctx context.Context, job types.DownloadJob, node pipeline.Node, comm chan<- pipeline.PipelineMessage) (types.DownloadJob, error) {
	client := http.Client{}
	err := utils.Download(ctx, &client, &job, make(map[string]string))
	return job, err
}
