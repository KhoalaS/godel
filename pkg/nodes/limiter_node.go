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

func LimiterNodeFunc(job types.DownloadJob, node pipeline.Node) (types.DownloadJob, error) {
	if limit, ex := node.Config["limit"]; ex {
		job.Limit = (limit).(int)
		return job, nil
	} else {
		return job, nil
	}
}
