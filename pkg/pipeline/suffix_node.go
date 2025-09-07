package pipeline

import "context"

func CreateSuffixNode() Node {
	return Node{
		Type:     "suffix",
		Name:     "Suffix",
		Category: NodeCategoryUtility,
		Io: map[string]*NodeIO{
			"input": {
				Id:        "input",
				ValueType: ValueTypeString,
				Required:  true,
				Label:     "Input",
				Type:      IOTypeConnectedOnly,
				Hooks: map[string]string{
					"output": "suffix",
				},
				HookMapping: map[string]string{
					"output": "input",
				},
			},
			"suffix": {
				Id:        "suffix",
				ValueType: ValueTypeString,
				Required:  true,
				Label:     "Suffix",
				Type:      IOTypeInput,
				Hooks: map[string]string{
					"output": "suffix",
				},
				HookMapping: map[string]string{
					"output": "suffix",
				},
			},
			"output": {
				Id:        "output",
				ValueType: ValueTypeString,
				Label:     "Output",
				Type:      IOTypeGenerated,
			},
		},
		Status: StatusPending,
		Run:    SuffixNodeFunc,
	}
}

func SuffixNodeFunc(ctx context.Context, node Node, pipelineId string, nodeId string) error {
	return nil
}
