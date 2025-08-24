package nodes

import (
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
)

func CreateLimiterNode() pipeline.Node {
	inputs := []pipeline.NodeInput{
		{
			Id:    "limit",
			Type:  "number",
			Label: "Limit",
		},
	}
	return pipeline.Node{
		Id:     "limiter",
		Phase:  pipeline.PrePhase,
		Run:    LimiterNodeFunc,
		Inputs: inputs,
		Name:   "Limiter",
		Status: pipeline.StatusPending,
	}
}

func LimiterNodeFunc(job *types.DownloadJob, data map[string]any) (types.DownloadJob, error) {
	next := *job.Clone()
	if limit, ex := data["limit"]; ex {
		next.Limit = (limit).(int)
		return next, nil
	} else {
		return next, nil
	}
}
