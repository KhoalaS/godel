package pipeline

import "fmt"

type InvalidNodeIOError struct {
	nodeType string
	nodeName string
	ioName   string
}

func NewInvalidNodeIOError(node *Node, ioName string) *InvalidNodeIOError {
	return &InvalidNodeIOError{
		nodeType: node.Type,
		nodeName: node.Name,
		ioName:   ioName,
	}
}

func (e *InvalidNodeIOError) Error() string {
	return fmt.Sprintf("[InvalidNodeIOError]: Invalid value at node '%s' with id '%s'. The IO value for '%s' was invalid.", e.nodeName, e.nodeType, e.ioName)
}
