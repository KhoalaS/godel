package nodes

import (
	"context"
	"net/http"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func CreateDownloaderNode() pipeline.Node {
	return pipeline.Node{
		Type: "default_downloader",
		Run:  DownloaderNodeFunc,
		Inputs: map[string]*pipeline.NodeIO{
			"limit": {
				Id:       "limit",
				Type:     pipeline.IOTypeNumber,
				Label:    "Limit",
				Required: false,
			},
			"url": {
				Id:       "url",
				Type:     pipeline.IOTypeString,
				Label:    "Url",
				Required: true,
			},
			"output_dir": {
				Id:       "output_dir",
				Type:     pipeline.IOTypeDirectory,
				Label:    "Output directory",
				Required: true,
			},
			"filename": {
				Id:       "filename",
				Type:     pipeline.IOTypeString,
				Label:    "Filename",
				Required: true,
			},
		},
		Name:     "Standard Downloader",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.DownloadNode,
		Outputs:  make(map[string]*pipeline.NodeIO),
	}
}

func DownloaderNodeFunc(ctx context.Context, node pipeline.Node, comm chan<- pipeline.PipelineMessage) error {
	client := http.Client{}

	job := types.NewDownloadJob()
	if node.Inputs["limit"] != nil {
		job.Limit = (node.Inputs["limit"].Value).(int)
	}

	job.Url = (node.Inputs["url"].Value).(string)
	job.DestPath = (node.Inputs["output_dir"].Value).(string)
	job.Filename = (node.Inputs["filename"].Value).(string)

	err := utils.Download(ctx, &client, job, make(map[string]string))
	return err
}
