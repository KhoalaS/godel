package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func NewRealdebridNode() Node {
	return Node{
		Type: "rd-downloader",
		Run:  RealdebridNodeFunc,
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
		Name:     "Real-Debrid Downloader",
		Status:   StatusPending,
		Category: NodeCategoryDownloader,
	}
}

func RealdebridNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	job := types.NewDownloadJob()

	jobUrl, ok := utils.FromAny[string](node.Io["url"].Value).Value()
	if !ok || jobUrl == "" {
		return NewInvalidNodeIOError(&node, "url")
	}

	job.Url = jobUrl

	rdJob, err := RealDebridTransformer(job)
	if err != nil {
		return err
	}

	node.Io["job"].Value = &rdJob

	return nil
}
