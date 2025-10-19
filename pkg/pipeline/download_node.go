package pipeline

import (
	"context"
	"net/http"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func NewDownloadNode() Node {
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

	_job, ok := utils.FromAny[*types.DownloadJob](node.Io["job"].Value).Value()

	if !ok || _job == nil {
		return NewInvalidNodeIOError(&node, "job")
	}

	job := _job.Clone()

	limit, ok := utils.FromAny[int](node.Io["limit"].Value).Value()
	if ok {
		job.Limit = int(limit)
	}

	destPath, ok := utils.FromAny[string](node.Io["output_dir"].Value).Value()
	if ok && destPath != "" {
		job.DestPath = destPath
	}

	filename, ok := utils.FromAny[string](node.Io["filename"].Value).Value()
	if ok && filename != "" {
		job.Filename = filename
	}

	file, err := Download(ctx, &client, &job, pipeline, node.Id)
	node.Io["file"].Value = file
	return err
}
