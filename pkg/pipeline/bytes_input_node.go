package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/utils"
)

const (
	Byte = 1
	KB   = 1 << 10 // 1024 bytes
	MB   = 1 << 20 // 1,048,576 bytes
	GB   = 1 << 30 // 1,073,741,824 bytes
)

func NewBytesInputNode() Node {
	return Node{
		Type:     "bytes-input",
		Run:      BytesInputNodeFunc,
		Name:     "Bytes",
		Status:   StatusPending,
		Category: NodeCategoryInput,
		Io: map[string]*NodeIO{
			"bytes": {
				Type:      IOTypeGenerated,
				Id:        "bytes",
				ValueType: ValueTypeNumber,
				Label:     "Bytes",
			},
			"amount": {
				Type:      IOTypeInput,
				ValueType: ValueTypeNumber,
				Id:        "amount",
				Label:     "Amount",
				Required:  true,
				Value:     10,
				Hooks: map[string]string{
					"bytes": "toBytes",
				},
				HookMapping: map[string]string{
					"bytes": "amount",
				},
			},
			"unit": {
				Type:      IOTypeSelection,
				ValueType: ValueTypeString,
				Id:        "unit",
				Label:     "Unit",
				Required:  true,
				Value:     "MB",
				Options:   []string{"B", "KB", "MB", "GB"},
				Hooks: map[string]string{
					"bytes": "toBytes",
				},
				HookMapping: map[string]string{
					"bytes": "unit",
				},
			},
		},
	}
}

func BytesInputNodeFunc(ctx context.Context, node Node, pipeline IPipeline) error {
	amount, ok := utils.FromAny[float64](node.Io["amount"].Value).Value()
	if !ok {
		return NewInvalidNodeIOError(&node, "amount")
	}

	bytesValue := 0
	unit, ok := utils.FromAny[string](node.Io["unit"].Value).Value()
	if !ok {
		return NewInvalidNodeIOError(&node, "unit")
	}

	switch unit {
	case "B":
		bytesValue = int(amount)
	case "KB":
		bytesValue = int(amount) * KB
	case "MB":
		bytesValue = int(amount) * MB
	case "GB":
		bytesValue = int(amount) * GB
	default:
		return NewInvalidNodeIOError(&node, "unit")
	}

	node.Io["bytes"].Value = bytesValue

	return nil
}
