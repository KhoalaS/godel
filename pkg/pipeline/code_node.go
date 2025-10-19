package pipeline

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/KhoalaS/godel/pkg/custom_error"
	"github.com/KhoalaS/godel/pkg/runner"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/rs/zerolog/log"
)

var jsInputRegex = regexp.MustCompile(`(?s)function\s([a-zA-Z]+?)\s*\([a-zA-Z,\s\t]*?\)\s*\{.*?\}`)

func NewCodeNode() Node {
	return Node{
		Type:     "code",
		Run:      CodeNodeFunc,
		Name:     "Code",
		Status:   StatusPending,
		Category: NodeCategoryInput,
		Io: map[string]*NodeIO{
			"input": {
				Type:      IOTypeInput,
				Id:        "input",
				ValueType: ValueTypeCode,
				Label:     "JS Code",
			},
			"output": {
				Type:      IOTypeGenerated,
				ValueType: ValueTypeUnknown,
				Id:        "output",
				Label:     "Output",
			},
		},
	}
}

func CodeNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {

	var r runner.PorfforRunner

	input, ok := utils.FromAny[string](node.Io["input"].Value).Value()
	if !ok {
		return NewInvalidNodeIOError(&node, "input")
	}

	m := jsInputRegex.FindStringSubmatch(input)
	if len(m) != 2 {
		return custom_error.FromError(errors.New("input is not a valid JS function"), 2, "CodeNodeFunc")
	}

	name := m[1]

	wrappedInput := fmt.Sprintf("%s;console.log(%s())", input, name)

	result, err := r.Run(wrappedInput)
	if err != nil {
		return err
	}

	node.Io["output"].Value = result

	log.Debug().Any("result", result).Send()

	return nil
}
