package pipeline

import (
	"context"
)

func NewBasenameNode() Node {
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
				HookMapping: map[string]string{
					"basename": "path",
				},
			},
			"basename": {
				Type:      IOTypeGenerated,
				Id:        "basename",
				ValueType: ValueTypeString,
				Label:     "Basename",
			},
		},
		Name:     "Basename",
		Status:   StatusPending,
		Category: NodeCategoryUtility,
	}
}

func BasenameNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	return nil
}
