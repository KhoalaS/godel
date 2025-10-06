package pipeline

import (
	"context"
	"errors"
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
			"file": {
				Type:      IOTypeGenerated,
				Id:        "file",
				ValueType: ValueTypeFile,
				Label:     "File",
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

	_job, ok := (node.Io["job"].Value).(*types.DownloadJob)

	if !ok || _job == nil {
		return errors.New("invalid download job input")
	}

	job := _job.Clone()

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

	file, err := Download(ctx, &client, &job, pipeline, node.Id)
	node.Io["file"].Value = file
	return err
}
