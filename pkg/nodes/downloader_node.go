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
		Type: "default-downloader",
		Run:  DownloaderNodeFunc,
		Io: map[string]*pipeline.NodeIO{
			"limit": {
				Type:      pipeline.IOTypeInput,
				Id:        "limit",
				ValueType: pipeline.ValueTypeNumber,
				Label:     "Limit",
				Required:  false,
			},
			"url": {
				Type:      pipeline.IOTypeInput,
				Id:        "url",
				ValueType: pipeline.ValueTypeString,
				Label:     "Url",
				Required:  true,
			},
			"output_dir": {
				Type:      pipeline.IOTypeInput,
				Id:        "output_dir",
				ValueType: pipeline.ValueTypeDirectory,
				Label:     "Output directory",
				Required:  true,
			},
			"filename": {
				Type:      pipeline.IOTypeInput,
				Id:        "filename",
				ValueType: pipeline.ValueTypeString,
				Label:     "Filename",
				Required:  true,
			},
		},
		Name:     "Standard Downloader",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.DownloaderNode,
	}
}

func DownloaderNodeFunc(ctx context.Context, node pipeline.Node, comm chan<- pipeline.PipelineMessage) error {
	client := http.Client{}

	job := types.NewDownloadJob()
	if node.Io["limit"] != nil {
		job.Limit = (node.Io["limit"].Value).(int)
	}

	job.Url = (node.Io["url"].Value).(string)
	job.DestPath = (node.Io["output_dir"].Value).(string)
	job.Filename = (node.Io["filename"].Value).(string)

	err := utils.Download(ctx, &client, job, make(map[string]string))
	return err
}
