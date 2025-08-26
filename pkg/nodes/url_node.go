package nodes

import (
	"context"
	"errors"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
)

func CreateUrlNode() pipeline.Node {
	inputs := []pipeline.NodeInput{
		{
			Id:    "url",
			Type:  "string",
			Label: "Url",
		},
	}
	return pipeline.Node{
		Type:   "url",
		Phase:  pipeline.PrePhase,
		Run:    UrlNodeFunc,
		Inputs: inputs,
		Name:   "Url",
		Status: pipeline.StatusPending,
	}
}

func UrlNodeFunc(ctx context.Context, job types.DownloadJob, node pipeline.Node, comm chan<- pipeline.PipelineMessage) (types.DownloadJob, error) {
	if _url, ex := node.Config["url"]; ex {
		var ok bool
		job.Url, ok = (_url).(string)
		if !ok {
			return job, errors.New("url was not a string")
		}
		return job, nil
	} else {
		return job, nil
	}
}
