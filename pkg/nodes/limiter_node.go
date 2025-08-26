package nodes

import (
	"context"
	"errors"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
)

func CreateLimiterNode() pipeline.Node {
	inputs := map[string]pipeline.NodeIO{
		"limit": {
			Id:       "limit",
			Type:     "number",
			Label:    "Limit",
			Required: true,
			Value:    1000000,
		},
	}
	return pipeline.Node{
		Type:     "limiter",
		Phase:    pipeline.PrePhase,
		Run:      LimiterNodeFunc,
		Inputs:   inputs,
		Name:     "Limiter",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.ConnectorNode,
	}
}

func LimiterNodeFunc(ctx context.Context, job types.DownloadJob, node pipeline.Node, comm chan<- pipeline.PipelineMessage) (types.DownloadJob, error) {
	if limit, ex := node.Inputs["limit"]; ex {
		var ok bool
		job.Limit, ok = (limit.Value).(int)
		if !ok {
			return job, errors.New("limit was not a number")
		}
		return job, nil
	} else {
		return job, nil
	}
}
