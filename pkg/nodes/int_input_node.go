package nodes

import (
	"context"

	"github.com/KhoalaS/godel/pkg/pipeline"
)

func CreateIntInputNode() pipeline.Node {
	return pipeline.Node{
		Type:     "int_input",
		Run:      IntInputNodeFunc,
		Inputs:   make(map[string]*pipeline.NodeIO),
		Name:     "Integer Input",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.SourceNode,
		Outputs: map[string]*pipeline.NodeIO{
			"output": {
				Id:       "output",
				Type:     pipeline.IOTypeNumber,
				Label:    "Output",
				Required: true,
				Value:    1000000,
			},
		},
	}
}

func IntInputNodeFunc(ctx context.Context, node pipeline.Node, comm chan<- pipeline.PipelineMessage) error {
	return nil
}
