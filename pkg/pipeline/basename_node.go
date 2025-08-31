package pipeline

import (
	"context"
)

func CreateBasenameNode() Node {
	return Node{
		Type: "basename",
		Run:  BasenameNodeFunc,
		Io: map[string]*NodeIO{
			"input": {
				Type:      IOTypeInput,
				Id:        "input",
				ValueType: ValueTypeString,
				Label:     "String",
				Required:  true,
				Hooks: map[string]string{
					"basename": "basename",
				},
			},
			"basename": {
				Type:      IOTypeGenerated,
				Id:        "basename",
				ValueType: ValueTypeString,
				Label:     "Basename",
				Required:  true,
			},
		},
		Name:     "Basename",
		Status:   StatusPending,
		NodeType: InputNode,
	}
}

func BasenameNodeFunc(ctx context.Context, node Node, comm chan<- PipelineMessage) error {
	return nil
}
