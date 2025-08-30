package nodes

import (
	"context"

	"github.com/KhoalaS/godel/pkg/pipeline"
)

func CreateIntInputNode() pipeline.Node {
	return pipeline.Node{
		Type:     "int-input",
		Run:      IntInputNodeFunc,
		Name:     "Integer Input",
		Status:   pipeline.StatusPending,
		NodeType: pipeline.InputNode,
		Io: map[string]*pipeline.NodeIO{
			"output": {
				Type:      pipeline.IOTypePassthrough,
				Id:        "output",
				ValueType: pipeline.ValueTypeNumber,
				Label:     "Output",
				Required:  true,
				Value:     1000000,
			},
		},
	}
}

func IntInputNodeFunc(ctx context.Context, node pipeline.Node, comm chan<- pipeline.PipelineMessage) error {
	return nil
}
