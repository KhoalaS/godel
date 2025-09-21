package pipeline

import (
	"context"
	"strconv"
)

const (
	Byte = 1
	KB   = 1 << 10 // 1024 bytes
	MB   = 1 << 20 // 1,048,576 bytes
	GB   = 1 << 30 // 1,073,741,824 bytes
)

func CreateBytesInputNode() Node {
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
	b := 0
	amount := 0

	if node.Io["amount"] != nil && node.Io["amount"].Value != nil {
		switch v := node.Io["amount"].Value.(type) {
		case int:
			amount = v
		case float64:
			amount = int(v)
		case float32:
			amount = int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				amount = i
			}
		}
	}
	switch (node.Io["unit"].Value).(string) {
	case "B":
		b = amount
	case "KB":
		b = amount * KB
	case "MB":
		b = amount * MB
	case "GB":
		b = amount * GB
	}

	node.Io["bytes"].Value = b

	return nil
}
