package nodes

import (
	"context"
	"net/http"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func CreateDownloadNode() pipeline.Node {
	return pipeline.Node{
		Type: "download",
		Run:  DownloadNodeFunc,
		Io: map[string]*pipeline.NodeIO{
			"limit": {
				Type:      pipeline.IOTypeInput,
				Id:        "limit",
				ValueType: pipeline.ValueTypeNumber,
				Label:     "Limit",
				Required:  false,
			},
			"job": {
				Type:      pipeline.IOTypeInput,
				Id:        "job",
				ValueType: pipeline.ValueTypeDownloadJob,
				Label:     "Downloader",
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
				Type:      pipeline.IOTypePassthrough,
				Id:        "filename",
				ValueType: pipeline.ValueTypeString,
				Label:     "Filename",
				Required:  true,
			},
		},
		Name:     "Download",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.DownloaderNode,
	}
}

func DownloadNodeFunc(ctx context.Context, node pipeline.Node, comm chan<- pipeline.PipelineMessage) error {
	client := http.Client{}

	job := (node.Io["job"].Value).(*types.DownloadJob)

	if node.Io["limit"] != nil {
		job.Limit = (node.Io["limit"].Value).(int)
	}

	job.DestPath = (node.Io["output_dir"].Value).(string)
	job.Filename = (node.Io["filename"].Value).(string)

	err := utils.Download(ctx, &client, job, make(map[string]string))
	return err
}
