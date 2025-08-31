package pipeline

import (
	"context"
)

func CreateIntInputNode() Node {
	return Node{
		Type:     "int-input",
		Run:      IntInputNodeFunc,
		Name:     "Integer",
		Status:   StatusPending,
		NodeType: InputNode,
		Io: map[string]*NodeIO{
			"output": {
				Type:      IOTypePassthrough,
				Id:        "output",
				ValueType: ValueTypeNumber,
				Required:  true,
				Value:     1000000,
			},
		},
	}
}

func IntInputNodeFunc(ctx context.Context, node Node, comm chan<- PipelineMessage) error {
	return nil
}
