package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/types"
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
	job.Url = (node.Io["url"].Value).(string)

	rdJob, err := RealDebridTransformer(job)
	if err != nil {
		return err
	}

	node.Io["job"].Value = &rdJob

	return nil
}
