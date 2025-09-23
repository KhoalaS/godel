package pipeline

import (
	"context"
)

func CreateDisplayNode() Node {
	return Node{
		Type:     "display",
		Run:      DisplayNodeFunc,
		Name:     "Display",
		Status:   StatusPending,
		Category: NodeCategoryUtility,
		Io: map[string]*NodeIO{
			"input": {
				Type:      IOTypeInput,
				Id:        "input",
				ValueType: ValueTypeUnknown,
				Label:     "Input",
				Hooks: map[string]string{
					"display": "getValue",
				},
				HookMapping: map[string]string{
					"display": "input",
				},
			},
			"display": {
				Id:        "display",
				ValueType: ValueTypeUnknown,
				Type:      IOTypeOutput,
			},
		},
	}
}

func DisplayNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	return nil
}
