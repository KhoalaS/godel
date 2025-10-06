package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/KhoalaS/godel/pkg/utils"
)

func CreateUnrarNode() Node {
	return Node{
		Type:     "unrar",
		Run:      UnrarNodeFunc,
		Name:     "Unrar",
		Status:   StatusPending,
		Category: NodeCategoryUtility,
		Io: map[string]*NodeIO{
			"file": {
				Id:        "file",
				ValueType: ValueTypeFile,
				Label:     "File",
				Required:  true,
				Type:      IOTypeConnectedOnly,
			},
			"password": {
				Id:        "password",
				ValueType: ValueTypeString,
				Label:     "Password",
				Required:  false,
				Value:     "",
				Type:      IOTypeInput,
			},
		},
	}
}

func UnrarNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	unrarExists, _ := utils.ExecutableExists("unrar")

	if !unrarExists {
		// TODO unrar in go
		return nil
	}

	unrarCommand := exec.Command("unrar", "x")
	if password, ok := node.Io["password"].Value.(string); ok && password != "" {
		unrarCommand.Args = append(unrarCommand.Args, fmt.Sprintf("-p%s", password))
	}

	file, ok := node.Io["file"].Value.(IFile)
	if !ok {
		return errors.New("missing file input")
	}

	absolutePath, err := file.GetAbsolutePath()
	if err != nil {
		return err
	}

	unrarCommand.Args = append(unrarCommand.Args, absolutePath, file.GetDestinationFolder())
	err = unrarCommand.Start()
	if err != nil {
		return err
	}

	err = unrarCommand.Wait()
	if err != nil {
		return err
	}

	return nil
}
