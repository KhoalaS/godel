package pipeline

import (
	"context"
)

func CreateDirectoryInputNode() Node {
	return Node{
		Type:     "directory-input",
		Run:      DirectoryInputNodeFunc,
		Name:     "Directory",
		Status:   StatusPending,
		Category: NodeCategoryInput,
		Io: map[string]*NodeIO{
			"directory": {
				Type:      IOTypeOutput,
				Id:        "directory",
				ValueType: ValueTypeDirectory,
				Required:  true,
				Value:     "./",
			},
		},
	}
}

func DirectoryInputNodeFunc(ctx context.Context, node Node, comm chan<- PipelineMessage, pipelineId string, nodeId string) error {
	return nil
}
