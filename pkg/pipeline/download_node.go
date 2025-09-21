package pipeline

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KhoalaS/godel/pkg/types"
)

func CreateDownloadNode() Node {
	return Node{
		Type: "download",
		Run:  DownloadNodeFunc,
		Io: map[string]*NodeIO{
			"limit": {
				Type:      IOTypeInput,
				Id:        "limit",
				ValueType: ValueTypeNumber,
				Label:     "Limit (Bytes/s)",
				Required:  false,
			},
			"job": {
				Type:      IOTypeConnectedOnly,
				Id:        "job",
				ValueType: ValueTypeDownloadJob,
				Label:     "Downloader",
				Required:  true,
			},
			"output_dir": {
				Type:      IOTypeInput,
				Id:        "output_dir",
				ValueType: ValueTypeDirectory,
				Label:     "Output directory",
				Required:  true,
			},
			"filename": {
				Type:      IOTypePassthrough,
				Id:        "filename",
				ValueType: ValueTypeString,
				Label:     "Filename",
				Required:  true,
			},
		},
		Name:     "Download",
		Status:   StatusPending,
		Category: NodeCategoryUtility,
	}
}

func DownloadNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	client := http.Client{}

	job := (node.Io["job"].Value).(*types.DownloadJob).Clone()

	if node.Io["limit"] != nil && node.Io["limit"].Value != nil {
		switch v := node.Io["limit"].Value.(type) {
		case int:
			job.Limit = v
		case float64:
			job.Limit = int(v)
		case float32:
			job.Limit = int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				job.Limit = i
			}
		}
	}

	if job.DestPath == "" {
		job.DestPath = (node.Io["output_dir"].Value).(string)
	}

	if job.Filename == "" {
		job.Filename = (node.Io["filename"].Value).(string)
	}

	err := Download(ctx, &client, &job, pipeline, node.Id)
	return err
}
