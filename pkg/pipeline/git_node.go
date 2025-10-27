package pipeline

import (
	"context"
	"os"

	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/go-git/go-git/v6"
	"github.com/rs/zerolog/log"
)

func NewGitNode() Node {
	return Node{
		Type:     "git",
		Run:      GitNodeFunc,
		Name:     "Git",
		Category: NodeCategoryUtility,
		Status:   StatusPending,
		Io: map[string]*NodeIO{
			"repositoryUrl": {
				Id:        "repositoryUrl",
				ValueType: ValueTypeString,
				Label:     "Repository Url",
				Required:  true,
				Type:      IOTypePassthrough,
			},
			"destination": {
				Id:        "destination",
				ValueType: ValueTypeDirectory,
				Label:     "Destination",
				Required:  true,
				Value:     "./",
				Type:      IOTypeInput,
			},
			"command": {
				Id:        "command",
				ValueType: ValueTypeString,
				Label:     "Command",
				Options: []string{
					"clone", "push", "pull", "fetch", "rebase", "status",
				},
				Value: "clone",
				Type:  IOTypeSelection,
			},
		},
	}
}

func GitNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	repositoryUrl, ok := utils.FromAny[string](node.Io["repositoryUrl"].Value).Value()
	if repositoryUrl == "" || !ok {
		return NewInvalidNodeIOError(&node, "repositoryUrl")
	}

	command, ok := utils.FromAny[string](node.Io["command"].Value).Value()
	if command == "" || !ok {
		return NewInvalidNodeIOError(&node, "command")
	}

	destination, ok := utils.FromAny[string](node.Io["destination"].Value).Value()
	if !ok {
		return NewInvalidNodeIOError(&node, "destination")
	}
	if destination == "" {
		destination = "./"
	}

	switch command {
	case "clone":
		_, err := git.PlainCloneContext(ctx, destination, &git.CloneOptions{
			URL:      repositoryUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	default:
		log.Warn().Msg("unimplemented git command")
	}
	return nil
}
